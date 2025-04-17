package accounting

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
)

// --- Mock repos ---In-memory SQLite repos
type fakeAccountGroupRepo struct {
	groups []*AccountGroup
	err    error
}

func (f *fakeAccountGroupRepo) GetAll(ctx context.Context) ([]*AccountGroup, error) {
	return f.groups, f.err
}
func (f *fakeAccountGroupRepo) GetByName(ctx context.Context, name string) (AccountGroup, error) {
	return AccountGroup{}, nil
}

func (f *fakeAccountGroupRepo) Upsert(ctx context.Context, group *AccountGroup) error {
	return nil
}

func (r *fakeAccountGroupRepo) Insert(ctx context.Context, group *AccountGroup) error {
	return nil
}

type fakeAccountRepo struct {
	accounts []*Account
	err      error
}

func (f *fakeAccountRepo) GetAll(ctx context.Context) ([]*Account, error) {
	return f.accounts, f.err
}

func (f *fakeAccountRepo) ByName(ctx context.Context, name string) (Account, error) {
	return Account{}, nil
}

func (f *fakeAccountRepo) Save(ctx context.Context, account *Account) error {
	return nil
}

// --- Helpers for test convenience ---

// Recursively searches for a specific group within a ChartOfAccounts tree
func findNode(root *ChartOfAccountsNode, groupName string) *ChartOfAccountsNode {
	if root.Group != nil && root.Group.Name == groupName {
		return root
	}
	for _, child := range root.Children {
		if found := findNode(child, groupName); found != nil {
			return found
		}
	}

	return nil
}

// --- Test Cases ---

func TestBuildChartOfAccountsTree(t *testing.T) {
	ctx := context.Background()

	t.Run("Empty Repositories", func(t *testing.T) {
		groupRepo := &fakeAccountGroupRepo{groups: []*AccountGroup{}}
		accountRepo := &fakeAccountRepo{accounts: []*Account{}}

		tree, err := BuildChartOfAccountsTree(ctx, groupRepo, accountRepo)
		if err != nil {
			t.Fatalf("expected no error; got %v", err)
		}

		// virtual root should have no children
		if len(tree.Children) != 0 {
			t.Errorf("expected 0 children; got %d", len(tree.Children))
		}
	})

	t.Run("Error Propagation", func(t *testing.T) {
		expectedErr := errors.New("group repo error")
		groupRepo := &fakeAccountGroupRepo{err: expectedErr}
		accountRepo := &fakeAccountRepo{accounts: []*Account{}}

		_, err := BuildChartOfAccountsTree(ctx, groupRepo, accountRepo)
		if err == nil {
			t.Fatalf("expected error; got nil")
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("expected error %v; got %v", expectedErr, err)
		}

		// Now simulate error in account repo
		groupRepo = &fakeAccountGroupRepo{groups: []*AccountGroup{}}
		expectedErr = errors.New("account repo error")
		accountRepo = &fakeAccountRepo{err: expectedErr}

		_, err = BuildChartOfAccountsTree(ctx, groupRepo, accountRepo)
		if err == nil {
			t.Fatalf("expected error; got nil")
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("expected error %v; got %v", expectedErr, err)
		}
	})

	t.Run("Hierarchical Chart with Accounts", func(t *testing.T) {
		// Setup sample groups:
		// Assets (root)
		//   - Current Assets (child of Assets)
		// Liabilities (root)
		//   - Current Liabilities (child of Liabilities)
		assets := &AccountGroup{
			Name:         "Assets",
			ParentName:   sql.NullString{Valid: false},
			DisplayAfter: sql.NullString{Valid: false},
			IsImmutable:  false,
		}
		currentAssets := &AccountGroup{
			Name:         "Current Assets",
			ParentName:   sql.NullString{String: "Assets", Valid: true},
			DisplayAfter: sql.NullString{Valid: false},
			IsImmutable:  false,
		}
		liabilities := &AccountGroup{
			Name:         "Liabilities",
			ParentName:   sql.NullString{Valid: false},
			DisplayAfter: sql.NullString{Valid: false},
			IsImmutable:  false,
		}
		currentLiabilities := &AccountGroup{
			Name:         "Current Liabilities",
			ParentName:   sql.NullString{String: "Liabilities", Valid: true},
			DisplayAfter: sql.NullString{Valid: false},
			IsImmutable:  false,
		}

		groups := []*AccountGroup{liabilities, currentAssets, currentLiabilities, assets}

		// Setup sample accounts:
		// Cash and Inventory belong to "Current Assets"
		// Accounts Payable belongs to "Current Liabilities"
		cash := &Account{
			Name:            "Cash",
			ParentGroupName: "Current Assets",
			AccountType:     Asset,
			NormalBalance:   DebitNormal,
			DisplayAfter:    sql.NullString{Valid: false},
		}
		inventory := &Account{
			Name:            "Inventory",
			ParentGroupName: "Current Assets",
			AccountType:     Asset,
			NormalBalance:   DebitNormal,
			DisplayAfter:    sql.NullString{String: "Cash", Valid: true},
		}
		accountsPayable := &Account{
			Name:            "Accounts Payable",
			ParentGroupName: "Current Liabilities",
			AccountType:     Liability,
			NormalBalance:   CreditNormal,
			DisplayAfter:    sql.NullString{Valid: false},
		}
		accounts := []*Account{accountsPayable, inventory, cash}

		groupRepo := &fakeAccountGroupRepo{groups: groups}
		accountRepo := &fakeAccountRepo{accounts: accounts}

		tree, err := BuildChartOfAccountsTree(ctx, groupRepo, accountRepo)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Our virtual root should have two children: Assets and Liabilities.
		if len(tree.Children) != 2 {
			t.Fatalf("expected 2 root groups; got %d", len(tree.Children))
		}

		// Validate Assets branch
		assetsNode := findNode(tree, "Assets")
		if assetsNode == nil {
			t.Fatalf("Assets node not found")
		}
		currentAssetsNode := findNode(tree, "Current Assets")
		if currentAssetsNode == nil {
			t.Fatalf("Current Assets node not found")
		}

		// Check that "Current Assets" is a child of "Assets"
		found := false
		for _, child := range assetsNode.Children {
			if child == currentAssetsNode {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Current Assets node is not a child of Assets")
		}

		// Validate accounts attached to Current Assets
		if len(currentAssetsNode.Accounts) != 2 {
			t.Errorf("expected 2 accounts under Current Assets; got %d", len(currentAssetsNode.Accounts))
		} else {
			accountNames := []string{currentAssetsNode.Accounts[0].Name, currentAssetsNode.Accounts[1].Name}
			expectedNames := []string{"Cash", "Inventory"}
			// Use reflect.DeepEqual after sorting if order isn’t fixed
			if !reflect.DeepEqual(accountNames, expectedNames) {
				t.Errorf("expected accounts %v; got %v", expectedNames, accountNames)
			}
		}

		// Validate Liabilities branch
		liabilitiesNode := findNode(tree, "Liabilities")
		if liabilitiesNode == nil {
			t.Fatalf("Liabilities node not found")
		}
		currentLiabilitiesNode := findNode(tree, "Current Liabilities")
		if currentLiabilitiesNode == nil {
			t.Fatalf("Current Liabilities node not found")
		}
		found = false
		for _, child := range liabilitiesNode.Children {
			if child == currentLiabilitiesNode {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Current Liabilities node is not a child of Liabilities")
		}

		// Validate accounts attached to Current Liabilities
		if len(currentLiabilitiesNode.Accounts) != 1 {
			t.Errorf("expected 1 account under Current Liabilities; got %d", len(currentLiabilitiesNode.Accounts))
		} else if currentLiabilitiesNode.Accounts[0].Name != "Accounts Payable" {
			t.Errorf("expected account Accounts Payable; got %s", currentLiabilitiesNode.Accounts[0].Name)
		}

		// Check that accounts are in the correct order; Inventory is meant to come *after* Cash
		accountNames := []string{currentAssetsNode.Accounts[0].Name, currentAssetsNode.Accounts[1].Name}
		expectedNames := []string{"Cash", "Inventory"}
		if !reflect.DeepEqual(accountNames, expectedNames) {
			t.Errorf("expected accounts %v; got %v", expectedNames, accountNames)
		}
	})
}
