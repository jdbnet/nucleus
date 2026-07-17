package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"nucleus/internal/db"
	"nucleus/internal/models"
	"nucleus/internal/runner"
	"nucleus/internal/scheduler"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/targets", getTargets)
	mux.HandleFunc("POST /api/targets", createTarget)
	mux.HandleFunc("DELETE /api/targets/{id}", deleteTarget)
	mux.HandleFunc("POST /api/targets/{id}/scan", triggerScan)
	mux.HandleFunc("GET /api/scans", getScans)
	mux.HandleFunc("GET /api/scans/{id}/findings", getFindings)
}

func getTargets(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name, address, schedule, last_scan_at, created_at FROM targets ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var targets []models.Target
	for rows.Next() {
		var t models.Target
		rows.Scan(&t.ID, &t.Name, &t.Address, &t.Schedule, &t.LastScanAt, &t.CreatedAt)
		targets = append(targets, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(targets)
}

func createTarget(w http.ResponseWriter, r *http.Request) {
	var t models.Target
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := db.DB.Exec("INSERT INTO targets (name, address, schedule) VALUES (?, ?, ?)", t.Name, t.Address, t.Schedule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()
	t.ID = int(id)

	scheduler.ReloadScheduler()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func deleteTarget(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("DELETE FROM targets WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	scheduler.ReloadScheduler()
	w.WriteHeader(http.StatusNoContent)
}

func triggerScan(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Run async
	go runner.RunScan(id, true)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "scan triggered"}`))
}

func getScans(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT s.id, s.target_id, t.name, s.status, s.started_at, s.completed_at
		FROM scans s
		JOIN targets t ON s.target_id = t.id
		ORDER BY s.started_at DESC
		LIMIT 100
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var scans []models.Scan
	for rows.Next() {
		var s models.Scan
		rows.Scan(&s.ID, &s.TargetID, &s.TargetName, &s.Status, &s.StartedAt, &s.CompletedAt)
		
		// Query findings counts by severity
		counts := make(map[string]int)
		fRows, err := db.DB.Query("SELECT severity, COUNT(*) FROM findings WHERE scan_id = ? GROUP BY severity", s.ID)
		if err == nil {
			for fRows.Next() {
				var sev string
				var count int
				fRows.Scan(&sev, &count)
				counts[sev] = count
			}
			fRows.Close()
		}
		s.FindingCounts = counts

		scans = append(scans, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scans)
}

func getFindings(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query("SELECT id, scan_id, template_id, name, severity, host, matched_at, description, detected_at FROM findings WHERE scan_id = ? ORDER BY detected_at DESC", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var findings []models.Finding
	for rows.Next() {
		var f models.Finding
		rows.Scan(&f.ID, &f.ScanID, &f.TemplateID, &f.Name, &f.Severity, &f.Host, &f.MatchedAt, &f.Description, &f.DetectedAt)
		findings = append(findings, f)
	}
	
	if findings == nil {
		findings = []models.Finding{} // return empty array instead of null
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(findings)
}
