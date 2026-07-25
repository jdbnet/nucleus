package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"nucleus/internal/config"
)

var (
	currentSession string
	sessionMu      sync.Mutex
)

func generateSession() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := config.Get(config.KeyWebUser)
		pass := config.Get(config.KeyWebPass)
		if user == "" || pass == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("nucleus_session")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		sessionMu.Lock()
		valid := currentSession != "" && currentSession == cookie.Value
		sessionMu.Unlock()

		if !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	user := config.Get(config.KeyWebUser)
	pass := config.Get(config.KeyWebPass)
	if user == "" || pass == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if creds.Username != user || creds.Password != pass {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionMu.Lock()
	currentSession = generateSession()
	token := currentSession
	sessionMu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "nucleus_session",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionMu.Lock()
	currentSession = ""
	sessionMu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "nucleus_session",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	user := config.Get(config.KeyWebUser)
	pass := config.Get(config.KeyWebPass)
	if user == "" || pass == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"auth_required": false, "authenticated": true}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	cookie, err := r.Cookie("nucleus_session")
	if err != nil {
		w.Write([]byte(`{"auth_required": true, "authenticated": false}`))
		return
	}

	sessionMu.Lock()
	valid := currentSession != "" && currentSession == cookie.Value
	sessionMu.Unlock()

	if valid {
		w.Write([]byte(`{"auth_required": true, "authenticated": true}`))
	} else {
		w.Write([]byte(`{"auth_required": true, "authenticated": false}`))
	}
}
