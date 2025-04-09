package sqlite

import (
	// std

	"context"
	"database/sql"
	"slices"
	"testing"

	// internal
	"github.com/hoodnoah/ghoam/internal/accounting"
)

func TestAccountRepo_Save(t *testing.T) {

	t.Run("inserts a new account", func(t *testing.T) {
		ctx := context.Background()

		// Step 1: Create in-memory SQLite DB using our migration
		repos, err := New(":memory:")
		if err != nil {
			t.Fatalf("failed to create in-memory SQLite DB with error %v", err)
		}

		// Step 2: Define an account
		account := &accounting.Account{
			Name:            "Institution Bank X1234",
			ParentGroupName: "Assets",
			AccountType:     accounting.Asset,
			NormalBalance:   accounting.DebitNormal,
			DisplayAfter:    sql.NullString{},
		}

		// Step 3: save the account
		if err := repos.Accounts.Save(ctx, account); err != nil {
			t.Fatalf("failed to save account with error %v", err)
		}

		// Step 4: verify the account was saved
		var name string
		err = repos.Accounts.(*accountRepo).db.QueryRowContext(ctx, `SELECT name FROM accounts WHERE name = ?`, account.Name).Scan(&name)
		if err != nil {
			t.Fatalf("failed to query account with error %v", err)
		}

		if name != account.Name {
			t.Fatalf("expected account name %s, got %s", account.Name, name)
		}
	})

	t.Run("updates an existing account", func(t *testing.T) {
		ctx := context.Background()

		// Step 1: Create in-memory SQLite DB using our migration
		repos, err := New(":memory:")
		if err != nil {
			t.Fatalf("failed to create in-memory SQLite DB with error %v", err)
		}

		// Step 2: Define an account
		account := &accounting.Account{
			Name:            "Institution Bank X1234",
			ParentGroupName: "Assets",
			AccountType:     accounting.Asset,
			NormalBalance:   accounting.DebitNormal,
			DisplayAfter:    sql.NullString{},
		}

		// Step 3: save the account
		if err := repos.Accounts.Save(ctx, account); err != nil {
			t.Fatalf("failed to save account with error %v", err)
		}

		// Step 4: update the account (rename)
		account.AccountType = accounting.Liability
		if err := repos.Accounts.Save(ctx, account); err != nil {
			t.Fatalf("failed to update account with error %v", err)
		}

		// Step 5: Verify the account was updated
		var accountType string
		err = repos.Accounts.(*accountRepo).db.QueryRowContext(ctx, `SELECT account_type FROM accounts WHERE name = ?`, account.Name).Scan(&accountType)
		if err != nil {
			t.Fatalf("failed to query account with error %v", err)
		}
		if accountType != string(account.AccountType) {
			t.Fatalf("expected account type %s, got %s", account.AccountType, accountType)
		}
	})
}

func TestAccountRepo_GetAll(t *testing.T) {
	t.Run("retrieves all accounts", func(t *testing.T) {
		ctx := context.Background()

		// Step 1: Create in-memory SQLite DB using our migration
		repos, err := New(":memory:")
		if err != nil {
			t.Fatalf("failed to create in-memory SQLite DB with error %v", err)
		}

		// Step 2: Define a list of accounts
		accounts := []*accounting.Account{
			{
				Name:            "Institution Bank X1234",
				ParentGroupName: "Assets",
				AccountType:     accounting.Asset,
				NormalBalance:   accounting.DebitNormal,
				DisplayAfter:    sql.NullString{},
			},
			{
				Name:            "Accounts Receivable",
				ParentGroupName: "Assets",
				AccountType:     accounting.Asset,
				NormalBalance:   accounting.DebitNormal,
				DisplayAfter:    sql.NullString{String: "Institution Bank X1234", Valid: true},
			},
			{
				Name:            "Accounts Payable",
				ParentGroupName: "Liabilities",
				AccountType:     accounting.Liability,
				NormalBalance:   accounting.CreditNormal,
				DisplayAfter:    sql.NullString{String: "Accounts Receivable", Valid: true},
			},
		}

		// insert accounts
		for _, account := range accounts {
			if err := repos.Accounts.Save(ctx, account); err != nil {
				t.Fatalf("failed to save account %s with error %v", account.Name, err)
			}
		}

		// get accounts
		accountsRetrieved, err := repos.Accounts.GetAll(ctx)
		if err != nil {
			t.Fatalf("failed to get all accounts with error %v", err)
		}

		// Step 3: Verify the accounts were retrieved
		if len(accountsRetrieved) != len(accounts) {
			t.Fatalf("expected %d accounts, got %d", len(accounts), len(accountsRetrieved))
		}

		for _, account := range accounts {
			if !slices.ContainsFunc(accountsRetrieved, func(a *accounting.Account) bool {
				return account.Name == a.Name
			}) {
				t.Fatalf("expected account %s to be in the list", account.Name)
			}
		}

	})

	t.Run("retrieves all accounts in correct order", func(t *testing.T) {
		ctx := context.Background()

		// Step 1: Create in-memory SQLite DB using our migration
		repos, err := New(":memory:")
		if err != nil {
			t.Fatalf("failed to create in-memory SQLite DB with error %v", err)
		}

		// Step 2: Define a list of accounts
		accounts := []*accounting.Account{
			{
				Name:            "Institution Bank X1234",
				ParentGroupName: "Assets",
				AccountType:     accounting.Asset,
				NormalBalance:   accounting.DebitNormal,
				DisplayAfter:    sql.NullString{},
			},
			{
				Name:            "Accounts Receivable",
				ParentGroupName: "Assets",
				AccountType:     accounting.Asset,
				NormalBalance:   accounting.DebitNormal,
				DisplayAfter:    sql.NullString{String: "Institution Bank X1234", Valid: true},
			},
			{
				Name:            "Accounts Payable",
				ParentGroupName: "Liabilities",
				AccountType:     accounting.Liability,
				NormalBalance:   accounting.CreditNormal,
				DisplayAfter:    sql.NullString{String: "Accounts Receivable", Valid: true},
			},
		}

		// insert accounts
		for _, account := range accounts {
			if err := repos.Accounts.Save(ctx, account); err != nil {
				t.Fatalf("failed to save account %s with error %v", account.Name, err)
			}
		}

		// get accounts
		accountsRetrieved, err := repos.Accounts.GetAll(ctx)
		if err != nil {
			t.Fatalf("failed to get all accounts with error %v", err)
		}

		// Step 3: Verify the accounts were retrieved
		if len(accountsRetrieved) != len(accounts) {
			t.Fatalf("expected %d accounts, got %d", len(accounts), len(accountsRetrieved))
		}

		for i, account := range accounts {
			if accountsRetrieved[i].Name != account.Name {
				t.Fatalf("expected account %s at index %d, got %s", account.Name, i, accountsRetrieved[i].Name)
			}
		}
	})
}

func TestAccountRepo_GetByName(t *testing.T) {
	t.Run("retrieves a single account by name", func(t *testing.T) {
		ctx := context.Background()

		// Step 1: Create in-memory SQLite DB using our migration
		repos, err := New(":memory:")
		if err != nil {
			t.Fatalf("failed to create in-memory SQLite DB with error %v", err)
		}

		// Step 2: Define an account
		account := &accounting.Account{
			Name:            "Institution Bank X1234",
			ParentGroupName: "Assets",
			AccountType:     accounting.Asset,
			NormalBalance:   accounting.DebitNormal,
			DisplayAfter:    sql.NullString{},
		}

		// Step 3: save the account
		if err := repos.Accounts.Save(ctx, account); err != nil {
			t.Fatalf("failed to save account with error %v", err)
		}

		// Step 4: get the account by name
		accountRetrieved, err := repos.Accounts.ByName(ctx, account.Name)
		if err != nil {
			t.Fatalf("failed to get account by name with error %v", err)
		}

		if accountRetrieved.Name != account.Name {
			t.Fatalf("expected account name %s, got %s", account.Name, accountRetrieved.Name)
		}
	})

	t.Run("returns a specific error if the account is not found", func(t *testing.T) {
		ctx := context.Background()

		// Step 1: Create in-memory SQLite DB using our migration
		repos, err := New(":memory:")
		if err != nil {
			t.Fatalf("failed to create in-memory SQLite DB with error %v", err)
		}

		// Step 2: try and find an account
		_, err = repos.Accounts.ByName(ctx, "Nonexistent Account")

		if !accounting.IsAccountNotFound(err) {
			t.Fatalf("expected account not found error, got %v", err)
		}
	})
}
