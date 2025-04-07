package db

import (
	"database/sql"

	"github.com/hoodnoah/ghoam/internal/accounting"
)

func createSchema(db *sql.DB) error {
	// Foreign keys are off in Sqlite by default, so enable them
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}

	// prep the table creation queries
	creationQueries := []string{
		// account groups table
		`CREATE TABLE IF NOT EXISTS account_groups (
		  id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			parent_id TEXT REFERENCES account_groups(id),
			display_after TEXT REFERENCES account_groups(id),
			is_immutable BOOLEAN NOT NULL DEFAULT false
		);`,

		// account types table
		`CREATE TABLE IF NOT EXISTS account_types (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			display_after TEXT REFERENCES account_types(id)
		);`,

		// accounts table
		`CREATE TABLE IF NOT EXISTS accounts (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			parent_group_id TEXT NOT NULL REFERENCES account_groups(id),
			account_type_id TEXT NOT NULL REFERENCES account_types(id),
			display_after TEXT REFERENCES accounts(id),
			normal_balance TEXT NOT NULL CHECK (normal_balance IN ('debit', 'credit'))
		);`,

		// journal entries table
		`CREATE TABLE IF NOT EXISTS journal_entries (
			id TEXT PRIMARY KEY,
			timestamp TEXT NOT NULL,
			description TEXT,
			cross_reference TEXT REFERENCES journal_entries(id) 
		);`,

		// journal lines table
		`CREATE TABLE IF NOT EXISTS journal_lines (
			id TEXT PRIMARY KEY,
			account_id TEXT NOT NULL REFERENCES accounts(id),
			amount REAL NOT NULL,
			side TEXT NOT NULL CHECK (side IN ('debit', 'credit')),
			journal_entry_id TEXT NOT NULL REFERENCES journal_entries(id)
		);`,
	}

	for _, query := range creationQueries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

// create the totality of account types,
// asset, contra-asset, liability, equity, revenue, and expense
func bootStrapAccountTypes(db *sql.DB) error {
	accountTypes := []struct {
		ID           string
		Name         string
		DisplayAfter sql.NullString
	}{
		{ID: "asset", Name: "Asset", DisplayAfter: sql.NullString{}},
		{ID: "contra_asset", Name: "Contra Asset", DisplayAfter: sql.NullString{String: "asset", Valid: true}},
		{ID: "liability", Name: "Liability", DisplayAfter: sql.NullString{String: "contra_asset", Valid: true}},
		{ID: "equity", Name: "Equity", DisplayAfter: sql.NullString{String: "liability", Valid: true}},
		{ID: "revenue", Name: "Revenue", DisplayAfter: sql.NullString{String: "equity", Valid: true}},
		{ID: "expense", Name: "Expense", DisplayAfter: sql.NullString{String: "revenue", Valid: true}},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback() // safe, even on success

	statement, err := tx.Prepare(`
		INSERT INTO account_types
		  (id, name, display_after)
		VALUES (?, ?, ?);
	`)
	if err != nil {
		return err
	}

	defer statement.Close()

	for _, accountType := range accountTypes {
		if _, err := statement.Exec(accountType.ID, accountType.Name, accountType.DisplayAfter); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// adds asset, liability, and equity account types
// which are present in all accounting systems
func bootStrapImmutableAccountGroups(db *sql.DB) error {
	// create the immutable account groups
	immutableGroups := []accounting.AccountGroup{
		{
			ID:           "assets",
			Name:         "Assets",
			ParentID:     sql.NullString{},
			DisplayAfter: sql.NullString{},
			IsImmutable:  true,
		},
		{
			ID:           "liabilities",
			Name:         "Liabilities",
			ParentID:     sql.NullString{},
			DisplayAfter: sql.NullString{String: "assets", Valid: true},
			IsImmutable:  true,
		},
		{
			ID:           "equity",
			Name:         "Equity",
			ParentID:     sql.NullString{},
			DisplayAfter: sql.NullString{String: "liabilities", Valid: true},
			IsImmutable:  true,
		},
		{
			ID:           "revenues",
			Name:         "Revenues",
			ParentID:     sql.NullString{},
			DisplayAfter: sql.NullString{},
			IsImmutable:  true,
		},
		{
			ID:           "expenses",
			Name:         "Expenses",
			ParentID:     sql.NullString{},
			DisplayAfter: sql.NullString{String: "revenues", Valid: true},
			IsImmutable:  true,
		},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // save, even on success

	statement, err := tx.Prepare(`
		INSERT INTO account_groups
		  (id, name, parent_id, display_after, is_immutable)
		VALUES (?, ?, ?, ?, TRUE);
	`)
	if err != nil {
		return err
	}

	defer statement.Close()

	for _, group := range immutableGroups {
		if _, err := statement.Exec(group.ID, group.Name, group.ParentID, group.DisplayAfter); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// BootstrapDB creates the database schema, and populates
// the database with the minimum set of entries
// common to all accounting systems
func BootstrapDB(db *sql.DB) error {
	if err := createSchema(db); err != nil {
		return err
	}

	// must precede the account groups, since
	// account groups references account types
	if err := bootStrapAccountTypes(db); err != nil {
		return err
	}

	if err := bootStrapImmutableAccountGroups(db); err != nil {
		return err
	}

	return nil
}
