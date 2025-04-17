package accounting

import "github.com/hoodnoah/ghoam/internal/ordering"

// Accounts ----------------------------------------
func SortAccounts(accs []Account) ([]Account, error) {
	return ordering.TopoSort(
		accs,
		func(a Account) string { return a.Name },
		func(a Account) (string, bool) {
			if a.DisplayAfter.Valid {
				return a.DisplayAfter.String, true
			}
			return "", false
		},
	)
}

func SortAccountsInPlace(accs []*Account) error {
	sorted, err := ordering.TopoSort(accs, func(a *Account) string { return a.Name }, func(a *Account) (string, bool) {
		if a.DisplayAfter.Valid {
			return a.DisplayAfter.String, true
		}
		return "", false
	})

	if err != nil {
		return err
	}

	copy(accs, sorted)
	return nil
}
