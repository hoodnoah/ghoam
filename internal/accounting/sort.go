package accounting

import "github.com/hoodnoah/ghoam/internal/ordering"

// Accounts ----------------------------------------
func SortAccounts(accs []Account) ([]Account, error) {
	return ordering.TopoSort(
		accs,
		func(a Account) string { return a.ID },
		func(a Account) (string, bool) {
			if a.DisplayAfter.Valid {
				return a.DisplayAfter.String, true
			}
			return "", false
		},
	)
}
