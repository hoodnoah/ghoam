package db_test

import (
	// std
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"testing"

	// module under test
	"github.com/hoodnoah/ghoam/internal/db"
)

// TestBootstrapDB tests the bootstrapping of the db
func TestBootstrapDB(t *testing.T) {
	testDB := NewTestDB(t)

	t.Run("bootstrapDB doesn't fail", func(t *testing.T) {
		if err := db.BootstrapDB(testDB); err != nil {
			t.Fatalf("failed to bootstrap db: %v", err)
		}
	})

	t.Run("bootstrapDB creates the correct tables", func(t *testing.T) {
		tables := []string{
			"account_groups",
			"account_types",
			"accounts",
			"journal_entries",
			"journal_lines",
		}

		for _, table := range tables {
			// check if the table exists
			query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`

			row := testDB.QueryRow(query, table)

			var name string
			if err := row.Scan(&name); err != nil {
				t.Fatalf("table %s does not exist: %v", table, err)
			}

			// check if the table has the correct columns
			columns := map[string]string{
				"account_groups":  "id, name, parent_id, display_after, is_immutable",
				"account_types":   "id, name, display_after",
				"accounts":        "id, name, parent_group_id, account_type_id, display_after, normal_balance",
				"journal_entries": "id, timestamp, description, cross_reference",
				"journal_lines":   "id, account_id, amount, side, journal_entry_id",
			}

			query = fmt.Sprintf(`PRAGMA table_info(%s)`, table)
			rows, err := testDB.Query(query)

			if err != nil {
				t.Fatalf("failed to query table info for %s: %v", table, err)
			}

			defer rows.Close()

			var columnNames []string
			for rows.Next() {
				var cid int
				var name string
				var typ string
				var notnull int
				var dfltValue sql.NullString
				var pk int

				if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
					t.Fatalf("failed to scan row: %v", err)
				}

				columnNames = append(columnNames, name)
			}

			if err := rows.Err(); err != nil {
				t.Fatalf("error iterating rows: %v", err)
			}

			// check if the columns match
			expectedColumns := columns[table]
			want := strings.Split(expectedColumns, ", ")
			for _, column := range columnNames {
				if !slices.Contains(want, column) {
					t.Fatalf("table %s has unexpected column %s", table, column)
				}
			}
		}
	})

	t.Run("bootstrapDB adds the expected default account types", func(t *testing.T) {
		// check if the default account types exist
		query := `SELECT name FROM account_types;`

		rows, err := testDB.Query(query)
		if err != nil {
			t.Fatalf("failed to query account_types: %v", err)
		}

		defer rows.Close()

		var accountTypes []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				t.Fatalf("failed to scan row: %v", err)
			}

			accountTypes = append(accountTypes, name)
		}

		if err := rows.Err(); err != nil {
			t.Fatalf("error iterating rows: %v", err)
		}

		expectedAccountTypes := []string{
			"Asset",
			"Contra Asset",
			"Liability",
			"Equity",
			"Revenue",
			"Expense",
		}

		for _, accountType := range accountTypes {
			if !slices.Contains(expectedAccountTypes, accountType) {
				t.Fatalf("unexpected account type %s", accountType)
			}
		}
	})

	t.Run("bootstrapDB adds the expected default account groups", func(t *testing.T) {
		query := `SELECT name FROM account_groups`

		rows, err := testDB.Query(query)
		if err != nil {
			t.Fatalf("failed to query account_groups: %v", err)
		}
		defer rows.Close()

		var accountGroups []string

		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				t.Fatalf("failed to scan row: %v", err)
			}

			accountGroups = append(accountGroups, name)
		}

		if err := rows.Err(); err != nil {
			t.Fatalf("error iterating rows: %v", err)
		}

		expectedAccountGroups := []string{
			"Assets",
			"Liabilities",
			"Equity",
			"Revenues",
			"Expenses",
		}

		for _, accountGroup := range accountGroups {
			if !slices.Contains(expectedAccountGroups, accountGroup) {
				t.Fatalf("unexpected account group %s", accountGroup)
			}
		}
	})
}
