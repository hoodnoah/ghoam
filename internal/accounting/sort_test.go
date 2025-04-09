package accounting

import (
	"database/sql"
	"testing"
)

func TestSortAccounts(t *testing.T) {
	unsorted := []Account{
		{ID: "accounts receivable", Name: "Accounts Receivable", ParentGroupID: "assets", AccountType: Asset, NormalBalance: DebitNormal, DisplayAfter: sql.NullString{String: "cash", Valid: true}},
		{ID: "cash", Name: "Cash", ParentGroupID: "assets", AccountType: Asset, NormalBalance: DebitNormal, DisplayAfter: sql.NullString{}},
		{ID: "fixed assets", Name: "Fixed Assets", ParentGroupID: "assets", AccountType: Asset, NormalBalance: DebitNormal, DisplayAfter: sql.NullString{String: "accounts receivable", Valid: true}},
	}

	got, err := SortAccounts(unsorted)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantOrder := []string{"cash", "accounts receivable", "fixed assets"}
	for i, a := range got {
		if a.ID != wantOrder[i] {
			t.Fatalf("expected %s at index %d, got %s", wantOrder[i], i, a.ID)
		}
	}
}
