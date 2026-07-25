package runner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"

	"nucleus/internal/db"
	"nucleus/internal/notifier"
	"nucleus/internal/models"
)

func RunScan(targetID int, _ bool) {
	// 1. Get Target
	row := db.DB.QueryRow("SELECT id, name, address, schedule FROM targets WHERE id = ?", targetID)
	var t models.Target
	err := row.Scan(&t.ID, &t.Name, &t.Address, &t.Schedule)
	if err != nil {
		log.Println("Failed to find target for scan:", err)
		return
	}

	// 2. Create Scan Record
	res, err := db.DB.Exec("INSERT INTO scans (target_id, status, started_at) VALUES (?, 'running', ?)", t.ID, time.Now())
	if err != nil {
		log.Println("Failed to create scan record:", err)
		return
	}
	scanID, _ := res.LastInsertId()

	// 3. Update target last_scan_at
	db.DB.Exec("UPDATE targets SET last_scan_at = ? WHERE id = ?", time.Now(), t.ID)

	// 4. Run Nuclei
	// -omit-raw keeps JSONL lines small; without it, request/response bodies
	// routinely exceed bufio.Scanner's default 64KiB token limit and findings are dropped.
	cmd := exec.Command("nuclei", "-target", t.Address, "-jsonl", "-silent", "-omit-raw")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Failed to get stdout pipe:", err)
		db.DB.Exec("UPDATE scans SET status = 'failed', completed_at = ? WHERE id = ?", time.Now(), scanID)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Println("Failed to start nuclei:", err)
		db.DB.Exec("UPDATE scans SET status = 'failed', completed_at = ? WHERE id = ?", time.Now(), scanID)
		return
	}

	var findings []models.Finding

	scanner := bufio.NewScanner(stdout)
	// Safety net for any remaining large lines (default MaxScanTokenSize is only 64KiB).
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		var nf models.NucleiFinding
		if err := json.Unmarshal(line, &nf); err != nil {
			log.Printf("Failed to parse nuclei finding: %v", err)
			continue
		}

		displayHost := nf.Host
		ipStr := displayHost
		if hostPart, _, err := net.SplitHostPort(displayHost); err == nil {
			ipStr = hostPart
		}
		if ip := net.ParseIP(ipStr); ip != nil {
			if names, err := net.LookupAddr(ip.String()); err == nil && len(names) > 0 {
				hostname := strings.TrimSuffix(names[0], ".")
				displayHost = fmt.Sprintf("%s (%s)", hostname, displayHost)
			}
		}

		severity := strings.ToLower(nf.Info.Severity)

		// Insert finding
		res, err := db.DB.Exec(`
			INSERT INTO findings (scan_id, template_id, name, severity, host, matched_at, description)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, scanID, nf.TemplateID, nf.Info.Name, severity, displayHost, nf.MatchedAt, nf.Info.Description)

		if err != nil {
			log.Printf("Failed to insert finding %s: %v", nf.TemplateID, err)
			continue
		}

		fid, _ := res.LastInsertId()
		findings = append(findings, models.Finding{
			ID:          int(fid),
			ScanID:      int(scanID),
			TemplateID:  nf.TemplateID,
			Name:        nf.Info.Name,
			Severity:    severity,
			Host:        displayHost,
			MatchedAt:   nf.MatchedAt,
			Description: nf.Info.Description,
		})
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading nuclei output for target %s: %v", t.Name, err)
	}

	err = cmd.Wait()
	status := "completed"
	if err != nil {
		log.Println("Nuclei finished with error:", err)
		status = "failed"
	}

	db.DB.Exec("UPDATE scans SET status = ?, completed_at = ? WHERE id = ?", status, time.Now(), scanID)
	log.Printf("Scan completed for target %s: status=%s findings=%d", t.Name, status, len(findings))

	notifier.NotifyScanComplete(t, int(scanID), findings, status)
}
