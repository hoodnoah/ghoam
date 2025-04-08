package sqlite

import (
	// std
	"context"
	"database/sql"
	"testing"

	// internal
	"github.com/hoodnoah/ghoam/internal/accounting"
)

func TestAccountRepo_Save(t *testing.T) {
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
}
