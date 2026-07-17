package runner

import (
	"bufio"
	"encoding/json"
	"log"
	"os/exec"
	"time"

	"nucleus/internal/db"
	"nucleus/internal/mailer"
	"nucleus/internal/models"
)

func RunScan(targetID int, isManual bool) {
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
	cmd := exec.Command("nuclei", "-target", t.Address, "-jsonl", "-silent")
	
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
	hasMediumOrHigher := false

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Bytes()
		var nf models.NucleiFinding
		if err := json.Unmarshal(line, &nf); err == nil {
			// Insert finding
			res, err := db.DB.Exec(`
				INSERT INTO findings (scan_id, template_id, name, severity, host, matched_at, description)
				VALUES (?, ?, ?, ?, ?, ?, ?)
			`, scanID, nf.TemplateID, nf.Info.Name, nf.Info.Severity, nf.Host, nf.MatchedAt, nf.Info.Description)
			
			if err == nil {
				fid, _ := res.LastInsertId()
				f := models.Finding{
					ID:          int(fid),
					ScanID:      int(scanID),
					TemplateID:  nf.TemplateID,
					Name:        nf.Info.Name,
					Severity:    nf.Info.Severity,
					Host:        nf.Host,
					MatchedAt:   nf.MatchedAt,
					Description: nf.Info.Description,
				}
				findings = append(findings, f)
				
				if f.Severity == "critical" || f.Severity == "high" || f.Severity == "medium" {
					hasMediumOrHigher = true
				}
			}
		}
	}

	err = cmd.Wait()
	status := "completed"
	if err != nil {
		log.Println("Nuclei finished with error:", err)
		status = "failed"
	}

	db.DB.Exec("UPDATE scans SET status = ?, completed_at = ? WHERE id = ?", status, time.Now(), scanID)

	// Send Email Report if manual or has medium+ severity
	if isManual || hasMediumOrHigher {
		mailer.SendReport(t, findings)
	}
}
