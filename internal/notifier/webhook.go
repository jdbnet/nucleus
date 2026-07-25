package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"nucleus/internal/db"
	"nucleus/internal/models"
)

func SendWebhooks(summary ScanSummary) {
	rows, err := db.DB.Query(`
		SELECT id, name, provider, url, enabled, options, created_at
		FROM webhooks
		WHERE enabled = 1
		ORDER BY id ASC
	`)
	if err != nil {
		log.Printf("Failed to load webhooks: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var wh models.Webhook
		var enabled int
		var createdAt *string
		if err := rows.Scan(&wh.ID, &wh.Name, &wh.Provider, &wh.URL, &enabled, &wh.Options, &createdAt); err != nil {
			log.Printf("Failed to scan webhook row: %v", err)
			continue
		}
		wh.Enabled = enabled == 1
		if createdAt != nil {
			wh.CreatedAt = *createdAt
		}
		sendWebhook(wh, summary)
	}
}

func SendWebhookTest(wh models.Webhook) error {
	summary := TestSummary()
	summary.ScanLink = ScanURL("", 0)
	return sendWebhookRequest(wh, summary, true)
}

func sendWebhook(wh models.Webhook, summary ScanSummary) {
	if err := sendWebhookRequest(wh, summary, false); err != nil {
		log.Printf("Failed to send webhook %q (%s): %v", wh.Name, wh.Provider, err)
		return
	}
	log.Printf("Successfully sent webhook %q (%s)", wh.Name, wh.Provider)
}

func sendWebhookRequest(wh models.Webhook, summary ScanSummary, isTest bool) error {
	body, contentType, err := buildWebhookPayload(wh, summary, isTest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, wh.URL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "Nucleus/1.0")

	if wh.Provider == "ntfy" {
		req.Header.Set("Title", summary.Title())
		if isTest {
			req.Header.Set("Title", "Nucleus Test Notification")
		}
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("webhook returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	return nil
}

func buildWebhookPayload(wh models.Webhook, summary ScanSummary, isTest bool) ([]byte, string, error) {
	text := summary.TextBody()
	if isTest {
		text = summary.TestTextBody()
	}

	switch strings.ToLower(wh.Provider) {
	case "discord":
		content := text
		if summary.ScanLink != "" && !isTest {
			content = summary.MarkdownBody()
		}
		payload := map[string]interface{}{
			"content": content,
		}
		if isTest {
			payload["content"] = text
		}
		b, err := json.Marshal(payload)
		return b, "application/json", err

	case "slack":
		payload := map[string]string{"text": text}
		b, err := json.Marshal(payload)
		return b, "application/json", err

	case "gotify":
		title := summary.Title()
		if isTest {
			title = "Nucleus Test Notification"
		}
		payload := map[string]interface{}{
			"title":    title,
			"message":  text,
			"priority": 5,
		}
		b, err := json.Marshal(payload)
		return b, "application/json", err

	case "ntfy":
		return []byte(text), "text/plain; charset=utf-8", nil

	case "teams":
		payload := map[string]interface{}{
			"@type":    "MessageCard",
			"@context": "http://schema.org/extensions",
			"summary":  summary.Title(),
			"themeColor": "0076D7",
			"title":    summary.Title(),
			"text":     strings.ReplaceAll(text, "\n", "\n\n"),
		}
		if isTest {
			payload["title"] = "Nucleus Test Notification"
			payload["summary"] = "Nucleus Test Notification"
		}
		b, err := json.Marshal(payload)
		return b, "application/json", err

	default:
		return []byte(text), "text/plain; charset=utf-8", nil
	}
}
