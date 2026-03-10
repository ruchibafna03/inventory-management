package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

// runMigrations applies the embedded SQL migration on startup.
// Safe to run multiple times — all statements use IF NOT EXISTS.
func runMigrations(db *sqlx.DB) {
	sql, err := os.ReadFile("migrations/001_initial.sql")
	if err != nil {
		// In production the file may not exist (baked into image differently).
		// Skip silently — tables were created when the image was first deployed.
		log.Printf("migrate: no migration file found, skipping (%v)", err)
		return
	}
	if _, err := db.Exec(string(sql)); err != nil {
		log.Fatalf("migrate: failed to apply migrations: %v", err)
	}
	log.Println("migrate: schema up to date")
}
