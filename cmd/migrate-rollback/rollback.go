package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/hoodnoah/ghoam/internal/persistence/sqlite/migrate"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
	// Set up command-line flags.
	steps := flag.Int("steps", -1, "Number of migration steps to roll back (negative value, e.g. -1 to roll back the last migration)")
	dbPath := flag.String("db", "./db.sqlite", "Path to SQLite database file")
	flag.Parse()

	// Open the database.
	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Optionally enable foreign keys.
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		log.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// Trigger the rollback.
	log.Printf("Rolling back %d migration step(s)...", *steps)
	if err := migrate.Rollback(db, *steps); err != nil {
		log.Fatalf("Rollback failed: %v", err)
	}

	log.Println("Rollback completed successfully.")
}
