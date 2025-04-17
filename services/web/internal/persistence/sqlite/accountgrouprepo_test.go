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

func TestAccountGroupRepo_Upsert(t *testing.T) {

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
		if err := repos.AccountGroups.Upsert(ctx, group); err != nil {
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

	t.Run("errors on trying to update an immutable account group", func(t *testing.T) {
		ctx := context.Background()

		repos, err := New(":memory:")
		if err != nil {
			t.Fatalf("failed to create in-memory SQLite DB with error %v", err)
		}

		newGroup := accounting.AccountGroup{
			Name:         "Assets",
			ParentName:   sql.NullString{},
			DisplayAfter: sql.NullString{},
			IsImmutable:  false,
		}

		// Try and commit the change
		err = repos.AccountGroups.Upsert(ctx, &newGroup)

		if !accounting.IsGroupImmutable(err) {
			t.Fatalf("expected to receive an IsGroupImmutable error, received %v", err)
		}
	})
}

func TestAccountGroupRepo_GetAll(t *testing.T) {
	t.Run("retrieves all extant account groups", func(t *testing.T) {
		ctx := context.Background()

		repos, err := New(":memory:")
		if err != nil {
			t.Fatalf("failed to create in memory SQLite DB with error %v", err)
		}

		expected := []accounting.AccountGroup{
			{
				Name:         "Assets",
				ParentName:   sql.NullString{},
				DisplayAfter: sql.NullString{},
				IsImmutable:  true,
			},
			{
				Name:         "Liabilities",
				ParentName:   sql.NullString{},
				DisplayAfter: sql.NullString{String: "Assets", Valid: true},
				IsImmutable:  true,
			},
			{
				Name:         "Equity",
				ParentName:   sql.NullString{},
				DisplayAfter: sql.NullString{String: "Liabilities", Valid: true},
				IsImmutable:  true,
			},
			{
				Name:         "Revenues",
				ParentName:   sql.NullString{String: "Equity", Valid: true},
				DisplayAfter: sql.NullString{},
				IsImmutable:  true,
			},
			{
				Name:         "Expenses",
				ParentName:   sql.NullString{String: "Equity", Valid: true},
				DisplayAfter: sql.NullString{String: "Revenues", Valid: true},
				IsImmutable:  true,
			},
		}

		actual, err := repos.AccountGroups.GetAll(ctx)
		if err != nil {
			t.Fatalf("failed to retrieve all account groups with error %v", err)
		}

		if len(actual) != len(expected) {
			t.Fatalf("expected to find %d AccountGroups, found %d", len(expected), len(actual))
		}

		for _, expectedGroup := range expected {
			if !slices.ContainsFunc(actual, func(ag *accounting.AccountGroup) bool {
				return ag.Name == expectedGroup.Name
			}) {
				t.Fatalf("expected to retrieve group with name %s, didn't", expectedGroup.Name)
			}
		}
	})

	t.Run("returns AccountGroups in correct order", func(t *testing.T) {
		ctx := context.Background()

		repos, err := New(":memory:")
		if err != nil {
			t.Fatalf("failed to create in-memory SQLite DB with error %v", err)
		}

		// add a higher-level account group
		newGroup := accounting.AccountGroup{
			Name:         "Cash",
			ParentName:   sql.NullString{String: "Assets", Valid: true},
			DisplayAfter: sql.NullString{},
			IsImmutable:  false,
		}
		if err := repos.AccountGroups.Upsert(ctx, &newGroup); err != nil {
			t.Fatalf("failed to add new group with error \"%s\"", err)
		}

		expected := []accounting.AccountGroup{
			{
				Name:         "Assets",
				ParentName:   sql.NullString{},
				DisplayAfter: sql.NullString{},
				IsImmutable:  true,
			},
			{
				Name:         "Cash",
				ParentName:   sql.NullString{String: "Assets", Valid: true},
				DisplayAfter: sql.NullString{},
				IsImmutable:  false,
			},
			{
				Name:         "Liabilities",
				ParentName:   sql.NullString{},
				DisplayAfter: sql.NullString{String: "Assets", Valid: true},
				IsImmutable:  true,
			},
			{
				Name:         "Equity",
				ParentName:   sql.NullString{},
				DisplayAfter: sql.NullString{String: "Liabilities", Valid: true},
				IsImmutable:  true,
			},
			{
				Name:         "Revenues",
				ParentName:   sql.NullString{String: "Equity", Valid: true},
				DisplayAfter: sql.NullString{},
				IsImmutable:  true,
			},
			{
				Name:         "Expenses",
				ParentName:   sql.NullString{String: "Equity", Valid: true},
				DisplayAfter: sql.NullString{String: "Revenues", Valid: true},
				IsImmutable:  true,
			},
		}

		// get all account groups
		actual, err := repos.AccountGroups.GetAll(ctx)
		if err != nil {
			t.Fatalf("failed to get all AccountGroups with error \"%s\"", err)
		}

		for i := range expected {
			if !actual[i].Equals(&expected[i]) {
				t.Fatalf("expected to see %s at index %d, received %s", expected[i].Name, i, actual[i].Name)
			}
		}
	})
}

func TestAccountGroupRepo_Insert(t *testing.T) {
	t.Run("fails to save an AccountGroup which already exists", func(t *testing.T) {
		ctx := context.Background()

		repos, err := New(":memory:")
		if err != nil {
			t.Fatalf("failed to create in memory SQLite DB with error %v", err)
		}

		testGroup := accounting.AccountGroup{
			Name:         "Current Assets",
			ParentName:   sql.NullString{String: "Assets", Valid: true},
			DisplayAfter: sql.NullString{},
			IsImmutable:  false,
		}

		// upsert a new AccountGroup
		err = repos.AccountGroups.Upsert(ctx, &testGroup)
		if err != nil {
			t.Fatalf("failed to upsert a new account group for test setup: %v", err)
		}

		// try and *insert* the same account group a second time.
		err = repos.AccountGroups.Insert(ctx, &testGroup)
		if err == nil {
			t.Fatal("expected an error on re-inserting an existing account group, received none")
		}
	})

	t.Run("fails to save an AccountGroup whose parent doesn't exist", func(t *testing.T) {
		ctx := context.Background()
		repos, err := New(":memory:")
		if err != nil {
			t.Fatalf("failed to create in memory SQLite DB with error %v", err)
		}

		testGroup := accounting.AccountGroup{
			Name:         "Notes Payable",
			ParentName:   sql.NullString{String: "Long-Term Liabilities", Valid: true},
			DisplayAfter: sql.NullString{},
			IsImmutable:  false,
		}

		// try and insert
		err = repos.AccountGroups.Insert(ctx, &testGroup)

		if !accounting.IsParentNameNotExists(err) {
			t.Fatalf("expected a ParentNameNotExists error, received %v", err)
		}
	})

	t.Run("fails to save an AccountGroup whose displayAfter doesn't exist", func(t *testing.T) {
		ctx := context.Background()
		repos, err := New(":memory:")
		if err != nil {
			t.Fatalf("failed to create an in memory SQLite DB with error %v", err)
		}

		testGroup := accounting.AccountGroup{
			Name:         "Notes Payable",
			ParentName:   sql.NullString{String: "Liabilities", Valid: true},
			DisplayAfter: sql.NullString{String: "Current Liabilities", Valid: true},
			IsImmutable:  false,
		}

		// try and insert
		err = repos.AccountGroups.Insert(ctx, &testGroup)
		if !accounting.IsDisplayAfterNotExists(err) {
			t.Fatalf("expected a DisplayAfterNotExists error, received %v", err)
		}
	})

	t.Run("inserts a well-formed AccountGroup successfully", func(t *testing.T) {
		ctx := context.Background()
		repos, err := New(":memory:")
		if err != nil {
			t.Fatalf("failed to create an in memory SQLite DB with error %v", err)
		}

		// insert a valid DisplayAfter
		err = repos.AccountGroups.Insert(ctx, &accounting.AccountGroup{
			Name:         "Current Liabilities",
			ParentName:   sql.NullString{String: "Liabilities", Valid: true},
			DisplayAfter: sql.NullString{},
			IsImmutable:  false,
		})
		if err != nil {
			t.Fatalf("failed to insert group \"Current Liabilities\" with error %v", err)
		}

		testGroup := accounting.AccountGroup{
			Name:         "Long Term Liabilities",
			ParentName:   sql.NullString{String: "Liabilities", Valid: true},
			DisplayAfter: sql.NullString{String: "Current Liabilities", Valid: true},
			IsImmutable:  false,
		}

		err = repos.AccountGroups.Insert(ctx, &testGroup)
		if err != nil {
			t.Fatalf("failed to insert group \"%s\" with error %v", testGroup.Name, err)
		}
	})
}
