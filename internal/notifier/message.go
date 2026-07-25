package notifier

import (
	"fmt"
	"sort"
	"strings"

	"nucleus/internal/models"
)

type ScanSummary struct {
	Target      models.Target
	ScanID      int
	Status      string
	ScanLink    string
	Findings    []models.Finding
	Counts      map[string]int
	TotalCount  int
}

func BuildSummary(target models.Target, scanID int, findings []models.Finding, status, scanLink string) ScanSummary {
	counts := map[string]int{"critical": 0, "high": 0, "medium": 0, "low": 0, "info": 0}
	for _, f := range findings {
		counts[strings.ToLower(f.Severity)]++
	}

	sorted := append([]models.Finding(nil), findings...)
	sort.Slice(sorted, func(i, j int) bool {
		return severityScore[strings.ToLower(sorted[i].Severity)] > severityScore[strings.ToLower(sorted[j].Severity)]
	})

	return ScanSummary{
		Target:     target,
		ScanID:     scanID,
		Status:     status,
		ScanLink:   scanLink,
		Findings:   sorted,
		Counts:     counts,
		TotalCount: len(findings),
	}
}

func (s ScanSummary) Title() string {
	return fmt.Sprintf("Nucleus Scan Report: %s", s.Target.Name)
}

func (s ScanSummary) TextBody() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", s.Title())
	fmt.Fprintf(&b, "Target: %s (%s)\n", s.Target.Name, s.Target.Address)
	fmt.Fprintf(&b, "Status: %s\n", s.Status)
	fmt.Fprintf(&b, "Findings: %d\n", s.TotalCount)

	if s.TotalCount > 0 {
		b.WriteString("\nSummary:\n")
		for _, sev := range []string{"critical", "high", "medium", "low", "info"} {
			if s.Counts[sev] > 0 {
				fmt.Fprintf(&b, "  %s: %d\n", sev, s.Counts[sev])
			}
		}

		b.WriteString("\nTop findings:\n")
		limit := 5
		if len(s.Findings) < limit {
			limit = len(s.Findings)
		}
		for i := 0; i < limit; i++ {
			f := s.Findings[i]
			fmt.Fprintf(&b, "  [%s] %s on %s\n", strings.ToUpper(f.Severity), f.Name, f.Host)
		}
		if len(s.Findings) > limit {
			fmt.Fprintf(&b, "  ... and %d more\n", len(s.Findings)-limit)
		}
	} else {
		b.WriteString("\nNo findings were detected.\n")
	}

	if s.ScanLink != "" {
		fmt.Fprintf(&b, "\nView scan: %s\n", s.ScanLink)
	}

	return b.String()
}

func (s ScanSummary) MarkdownBody() string {
	var b strings.Builder
	fmt.Fprintf(&b, "**%s**\n", s.Title())
	fmt.Fprintf(&b, "**Target:** %s (%s)\n", s.Target.Name, s.Target.Address)
	fmt.Fprintf(&b, "**Status:** %s\n", s.Status)
	fmt.Fprintf(&b, "**Findings:** %d\n", s.TotalCount)

	if s.TotalCount > 0 {
		b.WriteString("\n**Summary:**\n")
		for _, sev := range []string{"critical", "high", "medium", "low", "info"} {
			if s.Counts[sev] > 0 {
				fmt.Fprintf(&b, "- %s: %d\n", sev, s.Counts[sev])
			}
		}
	} else {
		b.WriteString("\nNo findings were detected.\n")
	}

	if s.ScanLink != "" {
		fmt.Fprintf(&b, "\n[View scan](%s)\n", s.ScanLink)
	}

	return b.String()
}

func TestSummary() ScanSummary {
	return ScanSummary{
		Target: models.Target{
			Name:    "Test Target",
			Address: "10.0.0.0/24",
		},
		ScanID:   0,
		Status:   "completed",
		ScanLink: "",
		Counts: map[string]int{
			"critical": 0,
			"high":     1,
			"medium":   2,
			"low":      0,
			"info":     0,
		},
		TotalCount: 3,
		Findings: []models.Finding{
			{Severity: "high", Name: "Example High Finding", Host: "10.0.0.1", TemplateID: "example-high"},
			{Severity: "medium", Name: "Example Medium Finding", Host: "10.0.0.2", TemplateID: "example-medium"},
		},
	}
}

func (s ScanSummary) TestTextBody() string {
	if s.ScanID == 0 && s.Target.Name == "Test Target" {
		var b strings.Builder
		b.WriteString("Nucleus Test Notification\n\n")
		b.WriteString("This is a test notification from Nucleus. If you received this, your webhook is configured correctly.\n")
		return b.String()
	}
	return s.TextBody()
}
