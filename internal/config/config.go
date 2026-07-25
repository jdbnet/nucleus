package config

import (
	"os"
	"strings"
	"sync"

	"nucleus/internal/db"
)

const (
	KeyInstanceURL       = "instance_url"
	KeySMTPHost          = "smtp_host"
	KeySMTPPort          = "smtp_port"
	KeySMTPFrom          = "smtp_from"
	KeySMTPTo            = "smtp_to"
	KeySMTPUser          = "smtp_user"
	KeySMTPPass          = "smtp_pass"
	KeyRetentionDays     = "retention_days"
	KeyNotifyMinSeverity = "notify_min_severity"
	KeyWebUser           = "web_user"
	KeyWebPass           = "web_pass"
)

const secretMask = "********"

var envFallback = map[string]string{
	KeyInstanceURL:       "NUCLEUS_URL",
	KeySMTPHost:          "SMTP_HOST",
	KeySMTPPort:          "SMTP_PORT",
	KeySMTPFrom:          "SMTP_FROM",
	KeySMTPTo:            "SMTP_TO",
	KeySMTPUser:          "SMTP_USER",
	KeySMTPPass:          "SMTP_PASS",
	KeyRetentionDays:     "RETENTION_DAYS",
	KeyWebUser:           "WEB_USER",
	KeyWebPass:           "WEB_PASS",
}

var defaults = map[string]string{
	KeyNotifyMinSeverity: "always",
	KeyRetentionDays:     "30",
}

var (
	mu    sync.RWMutex
	cache map[string]string
)

type Settings struct {
	InstanceURL       string `json:"instance_url"`
	SMTPHost          string `json:"smtp_host"`
	SMTPPort          string `json:"smtp_port"`
	SMTPFrom          string `json:"smtp_from"`
	SMTPTo            string `json:"smtp_to"`
	SMTPUser          string `json:"smtp_user"`
	SMTPPass          string `json:"smtp_pass"`
	RetentionDays     string `json:"retention_days"`
	NotifyMinSeverity string `json:"notify_min_severity"`
	WebUser           string `json:"web_user"`
	WebPass           string `json:"web_pass"`
}

func Init() {
	seedDefaults()
	Reload()
}

func seedDefaults() {
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM settings WHERE key = ?", KeyNotifyMinSeverity).Scan(&count)
	if err != nil || count > 0 {
		return
	}
	_, _ = db.DB.Exec("INSERT INTO settings (key, value) VALUES (?, ?)", KeyNotifyMinSeverity, "always")
}

func Reload() {
	rows, err := db.DB.Query("SELECT key, value FROM settings")
	if err != nil {
		return
	}
	defer rows.Close()

	values := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err == nil {
			values[key] = value
		}
	}

	mu.Lock()
	cache = values
	mu.Unlock()
}

func Get(key string) string {
	mu.RLock()
	value, ok := cache[key]
	mu.RUnlock()

	if ok && value != "" {
		return value
	}

	if envKey, ok := envFallback[key]; ok {
		if envVal := os.Getenv(envKey); envVal != "" {
			return envVal
		}
	}

	if def, ok := defaults[key]; ok {
		return def
	}

	return ""
}

func GetAll() Settings {
	s := Settings{
		InstanceURL:       Get(KeyInstanceURL),
		SMTPHost:          Get(KeySMTPHost),
		SMTPPort:          Get(KeySMTPPort),
		SMTPFrom:          Get(KeySMTPFrom),
		SMTPTo:            Get(KeySMTPTo),
		SMTPUser:          Get(KeySMTPUser),
		RetentionDays:     Get(KeyRetentionDays),
		NotifyMinSeverity: Get(KeyNotifyMinSeverity),
		WebUser:           Get(KeyWebUser),
	}

	if Get(KeySMTPPass) != "" {
		s.SMTPPass = secretMask
	}
	if Get(KeyWebPass) != "" {
		s.WebPass = secretMask
	}

	return s
}

func GetAllForAPI() Settings {
	return GetAll()
}

func Set(key, value string) error {
	_, err := db.DB.Exec(`
		INSERT INTO settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, key, value)
	if err != nil {
		return err
	}
	Reload()
	return nil
}

func SetMany(updates map[string]string) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for key, value := range updates {
		if _, err := stmt.Exec(key, value); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	Reload()
	return nil
}

func IsSecretMask(value string) bool {
	return value == secretMask
}

func NormalizeInstanceURL(url string) string {
	return strings.TrimRight(strings.TrimSpace(url), "/")
}
