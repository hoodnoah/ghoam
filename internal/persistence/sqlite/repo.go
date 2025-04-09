package sqlite

import (
	// std
	"context"
	"database/sql"

	// external
	"github.com/hoodnoah/ghoam/internal/persistence/sqlite/migrate"
	_ "github.com/mattn/go-sqlite3" // sqlite driver

	// internal
	"github.com/hoodnoah/ghoam/internal/accounting"
)

type accountRepo struct {
	db *sql.DB
}

// type journalEntryRepo struct {
// 	db *sql.DB
// }

type Repositories struct {
	Accounts accounting.AccountRepository
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
		Accounts: &accountRepo{db: db},
		// JournalEntries: &journalEntryRepo{db: db},
	}, nil
}

// accountRepo implementation

// Save inserts or updates an account in the database.
// More or less an upsert.
func (r *accountRepo) Save(ctx context.Context, account *accounting.Account) error {
	const query = `
		INSERT INTO accounts
			(id, name, parent_group_id, account_type_id, display_after, normal_balance)
	  VALUES
			(?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			parent_group_id = excluded.parent_group_id,
			account_type_id = excluded.account_type_id,
			display_after = excluded.display_after,
			normal_balance = excluded.normal_balance;
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		account.ID,
		account.Name,
		account.ParentGroupID,
		account.AccountType,
		account.DisplayAfter,
		account.NormalBalance,
	)

	return err
}

// Retrieves all accounts
func (r *accountRepo) GetAll(ctx context.Context) ([]*accounting.Account, error) {
	const query = `
		SELECT id, name, parent_group_id, account_type_id, display_after, normal_balance
		FROM accounts
		ORDER BY display_after, name;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*accounting.Account

	for rows.Next() {
		var account accounting.Account
		if err := rows.Scan(
			&account.ID,
			&account.Name,
			&account.ParentGroupID,
			&account.AccountType,
			&account.DisplayAfter,
			&account.NormalBalance,
		); err != nil {
			return nil, err
		}

		accounts = append(accounts, &account)
	}

	return accounts, nil
}
