package sqlite

import (
	// std
	"context"
	"database/sql"

	// external
	_ "github.com/mattn/go-sqlite3" // sqlite driver

	// internal
	"github.com/hoodnoah/ghoam/internal/accounting"
)

type accountRepo struct {
	db *sql.DB
}

// Save inserts or updates an account in the database.
// More or less an upsert.
func (r *accountRepo) Save(ctx context.Context, account *accounting.Account) error {
	const query = `
		INSERT INTO accounts
			(name, parent_group_name, account_type, display_after, normal_balance)
	  VALUES
			(?, ?, ?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			parent_group_name = excluded.parent_group_name,
			account_type = excluded.account_type,
			display_after = excluded.display_after,
			normal_balance = excluded.normal_balance;
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		account.Name,
		account.ParentGroupName,
		account.AccountType,
		account.DisplayAfter,
		account.NormalBalance,
	)

	return err
}

// Retrieves all accounts
func (r *accountRepo) GetAll(ctx context.Context) ([]*accounting.Account, error) {
	const query = `
		SELECT name, parent_group_name, account_type, display_after, normal_balance
		FROM accounts
		ORDER BY name;
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
			&account.Name,
			&account.ParentGroupName,
			&account.AccountType,
			&account.DisplayAfter,
			&account.NormalBalance,
		); err != nil {
			return nil, err
		}

		accounts = append(accounts, &account)
	}

	// sort accounts in-place
	if err := accounting.SortAccountsInPlace(accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

// Retrieves an account by name
//
// Returns ErrAccountNotFound if the account does not exist.
func (r *accountRepo) ByName(ctx context.Context, name string) (accounting.Account, error) {
	const query = `
		SELECT name, parent_group_name, account_type, display_after, normal_balance
		FROM accounts
		WHERE name = ?;
	`

	var account accounting.Account
	err := r.db.QueryRowContext(ctx, query, name).Scan(&account.Name, &account.ParentGroupName, &account.AccountType, &account.DisplayAfter, &account.NormalBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return accounting.Account{}, &accounting.ErrAccountNotFound{Name: name}
		}
		return accounting.Account{}, err
	}

	return account, nil
}
