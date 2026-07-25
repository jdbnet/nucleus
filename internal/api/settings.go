package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"nucleus/internal/config"
	"nucleus/internal/db"
	"nucleus/internal/models"
	"nucleus/internal/notifier"
)

func getSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config.GetAllForAPI())
}

func updateSettings(w http.ResponseWriter, r *http.Request) {
	var input config.Settings
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updates := map[string]string{}

	if input.InstanceURL != "" {
		parsed, err := url.Parse(strings.TrimSpace(input.InstanceURL))
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			http.Error(w, "instance_url must be a valid http or https URL", http.StatusBadRequest)
			return
		}
		updates[config.KeyInstanceURL] = config.NormalizeInstanceURL(input.InstanceURL)
	} else {
		updates[config.KeyInstanceURL] = ""
	}

	updates[config.KeySMTPHost] = strings.TrimSpace(input.SMTPHost)
	updates[config.KeySMTPPort] = strings.TrimSpace(input.SMTPPort)
	updates[config.KeySMTPFrom] = strings.TrimSpace(input.SMTPFrom)
	updates[config.KeySMTPTo] = strings.TrimSpace(input.SMTPTo)
	updates[config.KeySMTPUser] = strings.TrimSpace(input.SMTPUser)
	updates[config.KeyRetentionDays] = strings.TrimSpace(input.RetentionDays)
	updates[config.KeyNotifyMinSeverity] = strings.TrimSpace(input.NotifyMinSeverity)
	updates[config.KeyWebUser] = strings.TrimSpace(input.WebUser)

	if updates[config.KeyRetentionDays] == "" {
		updates[config.KeyRetentionDays] = "30"
	}
	if updates[config.KeyNotifyMinSeverity] == "" {
		updates[config.KeyNotifyMinSeverity] = "always"
	}
	if updates[config.KeyWebUser] == "" {
		updates[config.KeyWebPass] = ""
	}

	if input.SMTPPass != "" && !config.IsSecretMask(input.SMTPPass) {
		updates[config.KeySMTPPass] = input.SMTPPass
	}
	if input.WebPass != "" && !config.IsSecretMask(input.WebPass) {
		updates[config.KeyWebPass] = input.WebPass
	}

	if err := config.SetMany(updates); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config.GetAllForAPI())
}

func getWebhooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT id, name, provider, url, enabled, options, created_at
		FROM webhooks
		ORDER BY id ASC
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	webhooks := []models.Webhook{}
	for rows.Next() {
		wh, err := scanWebhook(rows)
		if err != nil {
			continue
		}
		webhooks = append(webhooks, wh)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhooks)
}

func createWebhook(w http.ResponseWriter, r *http.Request) {
	var wh models.Webhook
	if err := json.NewDecoder(r.Body).Decode(&wh); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateWebhook(wh); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	enabled := 0
	if wh.Enabled {
		enabled = 1
	}

	res, err := db.DB.Exec(`
		INSERT INTO webhooks (name, provider, url, enabled, options)
		VALUES (?, ?, ?, ?, ?)
	`, wh.Name, strings.ToLower(wh.Provider), strings.TrimSpace(wh.URL), enabled, nullIfEmpty(wh.Options))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	wh.ID = int(id)
	wh.Provider = strings.ToLower(wh.Provider)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wh)
}

func updateWebhook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var wh models.Webhook
	if err := json.NewDecoder(r.Body).Decode(&wh); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateWebhook(wh); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	enabled := 0
	if wh.Enabled {
		enabled = 1
	}

	_, err = db.DB.Exec(`
		UPDATE webhooks SET name = ?, provider = ?, url = ?, enabled = ?, options = ?
		WHERE id = ?
	`, wh.Name, strings.ToLower(wh.Provider), strings.TrimSpace(wh.URL), enabled, nullIfEmpty(wh.Options), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	wh.ID = id
	wh.Provider = strings.ToLower(wh.Provider)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wh)
}

func deleteWebhook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("DELETE FROM webhooks WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func testWebhook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	row := db.DB.QueryRow(`
		SELECT id, name, provider, url, enabled, options, created_at
		FROM webhooks WHERE id = ?
	`, id)

	wh, err := scanWebhookRow(row)
	if err != nil {
		http.Error(w, "webhook not found", http.StatusNotFound)
		return
	}

	if err := notifier.SendWebhookTest(wh); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "test notification sent"})
}

func validateWebhook(wh models.Webhook) error {
	if strings.TrimSpace(wh.Name) == "" {
		return errString("name is required")
	}
	if strings.TrimSpace(wh.URL) == "" {
		return errString("url is required")
	}
	parsed, err := url.Parse(strings.TrimSpace(wh.URL))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return errString("url must be a valid URL")
	}

	provider := strings.ToLower(strings.TrimSpace(wh.Provider))
	validProviders := map[string]bool{
		"discord": true, "slack": true, "gotify": true,
		"ntfy": true, "teams": true, "generic": true,
	}
	if !validProviders[provider] {
		return errString("invalid provider")
	}
	return nil
}

type errString string

func (e errString) Error() string { return string(e) }

func nullIfEmpty(value string) interface{} {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func scanWebhook(rows scannable) (models.Webhook, error) {
	return scanWebhookRow(rows)
}

func scanWebhookRow(row scannable) (models.Webhook, error) {
	var wh models.Webhook
	var enabled int
	var options *string
	var createdAt *string
	err := row.Scan(&wh.ID, &wh.Name, &wh.Provider, &wh.URL, &enabled, &options, &createdAt)
	if err != nil {
		return wh, err
	}
	wh.Enabled = enabled == 1
	if options != nil {
		wh.Options = *options
	}
	if createdAt != nil {
		wh.CreatedAt = *createdAt
	}
	return wh, nil
}
