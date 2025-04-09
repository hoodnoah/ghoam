package sqlite

import (
	// std

	"context"
	"database/sql"
	"testing"

	// internal
	"github.com/hoodnoah/ghoam/internal/accounting"
)

func TestAccountGroupRepo_Save(t *testing.T) {

	t.Run("inserts a new account group", func(t *testing.T) {
		ctx := context.Background()

		// Step 1: Create in-memory SQLite DB using our migration
		repos, err := New(":memory:")
		if err != nil {
			t.Fatalf("failed to create in-memory SQLite DB with error %v", err)
		}

		// Step 2: Define an account group
		group := &accounting.AccountGroup{
			Name:         "Cash",
			ParentName:   sql.NullString{String: "Assets", Valid: true},
			DisplayAfter: sql.NullString{},
			IsImmutable:  false,
		}

		// Step 3: save the account
		if err := repos.AccountGroups.Save(ctx, group); err != nil {
			t.Fatalf("failed to save account group with error \"%v\"", err)
		}

		// Step 4: verify the account group was saved
		extantGroup, err := repos.AccountGroups.GetByName(ctx, "Cash")
		if err != nil {
			t.Fatalf("expected to retrieve account group \"%s\", failed with error \"%v\"", group.Name, err)
		}

		if extantGroup.Name != group.Name {
			t.Fatalf("expected to retrieve a group named \"%s\", received \"%s\"", group.Name, extantGroup.Name)
		}
	})
}
