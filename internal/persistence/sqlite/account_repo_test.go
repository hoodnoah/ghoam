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
			ID:            "institution bank X1234",
			Name:          "Institution Bank X1234",
			ParentGroupID: "assets",
			AccountType:   accounting.Asset,
			NormalBalance: accounting.DebitNormal,
			DisplayAfter:  sql.NullString{},
		}

		// Step 3: save the account
		if err := repos.Accounts.Save(ctx, account); err != nil {
			t.Fatalf("failed to save account with error %v", err)
		}

		// Step 4: verify the account was saved
		var name string
		err = repos.Accounts.(*accountRepo).db.QueryRowContext(ctx, `SELECT name FROM accounts WHERE id = ?`, account.ID).Scan(&name)
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
			ID:            "institution bank X1234",
			Name:          "Institution Bank X1234",
			ParentGroupID: "assets",
			AccountType:   accounting.Asset,
			NormalBalance: accounting.DebitNormal,
			DisplayAfter:  sql.NullString{},
		}

		// Step 3: save the account
		if err := repos.Accounts.Save(ctx, account); err != nil {
			t.Fatalf("failed to save account with error %v", err)
		}

		// Step 4: update the account (rename)
		account.Name = "Institution Bank X1234 Sweep"
		if err := repos.Accounts.Save(ctx, account); err != nil {
			t.Fatalf("failed to update account with error %v", err)
		}

		// Step 5: Verify the account was updated
		var name string
		err = repos.Accounts.(*accountRepo).db.QueryRowContext(ctx, `SELECT name FROM accounts WHERE id = ?`, account.ID).Scan(&name)
		if err != nil {
			t.Fatalf("failed to query account with error %v", err)
		}
		if name != account.Name {
			t.Fatalf("expected account name %s, got %s", account.Name, name)
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
				ID:            "institution bank X1234",
				Name:          "Institution Bank X1234",
				ParentGroupID: "assets",
				AccountType:   accounting.Asset,
				NormalBalance: accounting.DebitNormal,
				DisplayAfter:  sql.NullString{},
			},
			{
				ID:            "accounts receivable",
				Name:          "Accounts Receivable",
				ParentGroupID: "assets",
				AccountType:   accounting.Asset,
				NormalBalance: accounting.DebitNormal,
				DisplayAfter:  sql.NullString{String: "institution bank X1234", Valid: true},
			},
			{
				ID:            "accounts payable",
				Name:          "Accounts Payable",
				ParentGroupID: "liabilities",
				AccountType:   accounting.Liability,
				NormalBalance: accounting.CreditNormal,
				DisplayAfter:  sql.NullString{String: "accounts receivable", Valid: true},
			},
		}

		// insert accounts
		for _, account := range accounts {
			if err := repos.Accounts.Save(ctx, account); err != nil {
				t.Fatalf("failed to save account %s with error %v", account.ID, err)
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
				return account.ID == a.ID
			}) {
				t.Fatalf("expected account %s to be in the list", account.ID)
			}
		}

	})

	// t.Run("retrieves all accounts in correct order", func(t *testing.T) {
	// 	ctx := context.Background()

	// 	// Step 1: Create in-memory SQLite DB using our migration
	// 	repos, err := New(":memory:")
	// 	if err != nil {
	// 		t.Fatalf("failed to create in-memory SQLite DB with error %v", err)
	// 	}

	// 	// Step 2: Define a list of accounts
	// 	accounts := []*accounting.Account{
	// 		{
	// 			ID:            "institution bank X1234",
	// 			Name:          "Institution Bank X1234",
	// 			ParentGroupID: "assets",
	// 			AccountType:   accounting.Asset,
	// 			NormalBalance: accounting.DebitNormal,
	// 			DisplayAfter:  sql.NullString{},
	// 		},
	// 		{
	// 			ID:            "accounts receivable",
	// 			Name:          "Accounts Receivable",
	// 			ParentGroupID: "assets",
	// 			AccountType:   accounting.Asset,
	// 			NormalBalance: accounting.DebitNormal,
	// 			DisplayAfter:  sql.NullString{String: "institution bank X1234", Valid: true},
	// 		},
	// 		{
	// 			ID:            "accounts payable",
	// 			Name:          "Accounts Payable",
	// 			ParentGroupID: "liabilities",
	// 			AccountType:   accounting.Liability,
	// 			NormalBalance: accounting.CreditNormal,
	// 			DisplayAfter:  sql.NullString{String: "accounts receivable", Valid: true},
	// 		},
	// 	}

	// 	// insert accounts
	// 	for _, account := range accounts {
	// 		if err := repos.Accounts.Save(ctx, account); err != nil {
	// 			t.Fatalf("failed to save account %s with error %v", account.ID, err)
	// 		}
	// 	}

	// 	// get accounts
	// 	accountsRetrieved, err := repos.Accounts.GetAll(ctx)
	// 	if err != nil {
	// 		t.Fatalf("failed to get all accounts with error %v", err)
	// 	}

	// 	// Step 3: Verify the accounts were retrieved
	// 	if len(accountsRetrieved) != len(accounts) {
	// 		t.Fatalf("expected %d accounts, got %d", len(accounts), len(accountsRetrieved))
	// 	}

	// 	for i, account := range accounts {
	// 		if accountsRetrieved[i].ID != account.ID {
	// 			t.Fatalf("expected account %s at index %d, got %s", account.ID, i, accountsRetrieved[i].ID)
	// 		}
	// 	}
	// })
}
