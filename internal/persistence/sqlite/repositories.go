package sqlite

import (
	// std

	"database/sql"

	// external
	"github.com/hoodnoah/ghoam/internal/persistence/sqlite/migrate"
	_ "github.com/mattn/go-sqlite3" // sqlite driver

	// internal
	"github.com/hoodnoah/ghoam/internal/accounting"
)

type Repositories struct {
	Accounts      accounting.AccountRepository
	AccountGroups accounting.AccountGroupRepository
	// JournalEntries accounting.JournalEntryRepository
}

// New opens/creates the DB, runs migrations, enables FK checks, and returns repositories
func New(path string) (*Repositories, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, err
	}

	if err := migrate.Up(db); err != nil {
		return nil, err
	}

	return &Repositories{
		Accounts:      &accountRepo{db: db},
		AccountGroups: &accountGroupRepo{db: db},
		// JournalEntries: &journalEntryRepo{db: db},
	}, nil
}
