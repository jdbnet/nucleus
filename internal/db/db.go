package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	_, err = DB.Exec(Schema)
	if err != nil {
		log.Fatal("Failed to create schema:", err)
	}
	
	// Enable foreign keys
	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Println("Warning: Failed to enable foreign keys:", err)
	}
	
	log.Println("Database initialized successfully at", filepath)
}
