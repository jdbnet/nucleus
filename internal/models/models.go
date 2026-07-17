package models

import "time"

type Target struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Address    string     `json:"address"`
	Schedule   string     `json:"schedule"`
	LastScanAt *time.Time `json:"last_scan_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

type Scan struct {
	ID          int        `json:"id"`
	TargetID    int        `json:"target_id"`
	TargetName  string     `json:"target_name"` // For UI
	Status      string     `json:"status"`      // "running", "completed", "failed"
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	FindingCounts map[string]int `json:"finding_counts,omitempty"` // aggregated by severity
}

type Finding struct {
	ID          int       `json:"id"`
	ScanID      int       `json:"scan_id"`
	TemplateID  string    `json:"template_id"`
	Name        string    `json:"name"`
	Severity    string    `json:"severity"`
	Host        string    `json:"host"`
	MatchedAt   string    `json:"matched_at"`
	Description string    `json:"description"`
	DetectedAt  time.Time `json:"detected_at"`
}

// Nuclei JSONL Output Structure (Partial for what we need)
type NucleiFinding struct {
	TemplateID string `json:"template-id"`
	Info       struct {
		Name        string `json:"name"`
		Severity    string `json:"severity"`
		Description string `json:"description"`
	} `json:"info"`
	Host      string `json:"host"`
	MatchedAt string `json:"matched-at"`
	Timestamp string `json:"timestamp"`
}
