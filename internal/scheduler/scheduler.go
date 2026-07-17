package scheduler

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/robfig/cron/v3"

	"nucleus/internal/db"
	"nucleus/internal/models"
	"nucleus/internal/runner"
)

var c *cron.Cron

func InitScheduler() {
	c = cron.New()
	
	// Automated Data Cleanup Job (Runs every day at midnight)
	c.AddFunc("@daily", func() {
		days := os.Getenv("RETENTION_DAYS")
		if days == "" {
			days = "30"
		}
		
		retentionDays, err := strconv.Atoi(days)
		if err != nil {
			log.Printf("Invalid RETENTION_DAYS environment variable: %s", days)
			retentionDays = 30
		}
		
		log.Printf("Running automated cleanup (Retention: %d days)", retentionDays)
		
		modifier := fmt.Sprintf("-%d days", retentionDays)
		res, err := db.DB.Exec("DELETE FROM scans WHERE started_at < date('now', ?)", modifier)
		if err != nil {
			log.Printf("Cleanup failed: %v", err)
			return
		}
		
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected > 0 {
			log.Printf("Cleanup complete: removed %d old scans.", rowsAffected)
		}
	})

	// Load all targets
	rows, err := db.DB.Query("SELECT id, name, address, schedule FROM targets")
	if err != nil {
		log.Fatal("Failed to load targets for scheduler:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Target
		if err := rows.Scan(&t.ID, &t.Name, &t.Address, &t.Schedule); err != nil {
			continue
		}
		
		if t.Schedule != "" && t.Schedule != "manual" {
			AddJob(t)
		}
	}

	c.Start()
	log.Println("Scheduler started successfully")
}

func AddJob(t models.Target) {
	_, err := c.AddFunc(t.Schedule, func() {
		log.Printf("Triggering scheduled scan for target %s (%s)", t.Name, t.Address)
		runner.RunScan(t.ID, false)
	})
	if err != nil {
		log.Printf("Failed to schedule target %s: %v", t.Name, err)
	}
}

func ReloadScheduler() {
	c.Stop()
	InitScheduler()
}
