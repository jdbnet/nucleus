package notifier

import (
	"fmt"
	"log"
	"strings"

	"nucleus/internal/config"
	"nucleus/internal/mailer"
	"nucleus/internal/models"
)

var severityScore = map[string]int{
	"critical": 5,
	"high":     4,
	"medium":   3,
	"low":      2,
	"info":     1,
}

func ShouldNotify(minSeverity string, findings []models.Finding) bool {
	minSeverity = strings.ToLower(strings.TrimSpace(minSeverity))
	if minSeverity == "" || minSeverity == "always" {
		return true
	}

	threshold, ok := severityScore[minSeverity]
	if !ok {
		return true
	}

	for _, f := range findings {
		if score, ok := severityScore[strings.ToLower(f.Severity)]; ok && score >= threshold {
			return true
		}
	}
	return false
}

func ScanURL(instanceURL string, scanID int) string {
	base := config.NormalizeInstanceURL(instanceURL)
	if base == "" || scanID <= 0 {
		return ""
	}
	return fmt.Sprintf("%s/scans/%d", base, scanID)
}

func NotifyScanComplete(target models.Target, scanID int, findings []models.Finding, status string) {
	minSeverity := config.Get(config.KeyNotifyMinSeverity)
	if !ShouldNotify(minSeverity, findings) {
		logSkip(target.Name, minSeverity, len(findings))
		return
	}

	instanceURL := config.Get(config.KeyInstanceURL)
	scanLink := ScanURL(instanceURL, scanID)
	summary := BuildSummary(target, scanID, findings, status, scanLink)

	mailer.SendReport(target, findings, status, scanLink)
	SendWebhooks(summary)
}

func logSkip(targetName, minSeverity string, findingCount int) {
	log.Printf("Skipping notifications for target %s: no findings at or above severity threshold %q (%d total findings)",
		targetName, minSeverity, findingCount)
}
