package accounting

import (
	"database/sql"
	"testing"
)

func TestSortAccounts(t *testing.T) {
	unsorted := []Account{
		{Name: "Accounts Receivable", ParentGroupName: "Assets", AccountType: Asset, NormalBalance: DebitNormal, DisplayAfter: sql.NullString{String: "Cash", Valid: true}},
		{Name: "Cash", ParentGroupName: "Assets", AccountType: Asset, NormalBalance: DebitNormal, DisplayAfter: sql.NullString{}},
		{Name: "Fixed Assets", ParentGroupName: "Assets", AccountType: Asset, NormalBalance: DebitNormal, DisplayAfter: sql.NullString{String: "Accounts Receivable", Valid: true}},
	}

	got, err := SortAccounts(unsorted)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantOrder := []string{"Cash", "Accounts Receivable", "Fixed Assets"}
	for i, a := range got {
		if a.Name != wantOrder[i] {
			t.Fatalf("expected %s at index %d, got %s", wantOrder[i], i, a.Name)
		}
	}
}
