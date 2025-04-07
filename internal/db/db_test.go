package db_test

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// New returns a *sql.DB that lives only for the duration of the test.
// t.Helper() marks the caller's line in case of failure.
func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// DSN explanation:
	//   file:memdb1           → any name is fine; needed for shared cache
	//   ?mode=memory          → keep it in RAM
	//   &cache=shared         → every connection sees the same schema/data
	//   &_fk=1                → turn ON foreign‑key checks automatically
	dsn := "file:memdb1?mode=memory&cache=shared&_fk=1"

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		t.Fatalf("open in‑memory db: %v", err)
	}

	// One connection is enough for most unit tests and avoids “database
	// is locked” surprises.
	db.SetMaxOpenConns(1)

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			log.Printf("close db: %v", err)
		}
	})
	return db
}
