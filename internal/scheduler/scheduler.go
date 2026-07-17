package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"

	"nucleus/internal/db"
	"nucleus/internal/models"
	"nucleus/internal/runner"
)

var c *cron.Cron

func InitScheduler() {
	c = cron.New()
	
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
