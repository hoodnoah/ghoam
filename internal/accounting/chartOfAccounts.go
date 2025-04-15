package accounting

import (
	"context"
	"sort"

	"github.com/hoodnoah/ghoam/internal/ordering"
)

type ChartOfAccountsNode struct {
	Group    *AccountGroup
	Children []*ChartOfAccountsNode
	Accounts []*Account
}

func BuildChartOfAccountsTree(ctx context.Context, groupRepo AccountGroupRepository, accountRepo AccountRepository) (*ChartOfAccountsNode, error) {
	// fetch all account groups
	groups, err := groupRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// fetch all accounts
	accounts, err := accountRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// build sorted tree of account groups
	sortedRoots, err := BuildAccountGroupTree(groups)
	if err != nil {
		return nil, err
	}

	// Build map: group name -> AccountGroupNode
	groupMap := make(map[string]*ChartOfAccountsNode)
	for _, root := range sortedRoots {
		populateGroupMap(root, groupMap)
	}

	// Attach accounts to groups
	for _, acct := range accounts {
		if node, ok := groupMap[acct.ParentGroupName]; ok {
			node.Accounts = append(node.Accounts, acct)
		}
	}

	// Build the tree
	var root ChartOfAccountsNode // virtual root node
	for _, node := range groupMap {
		if node.Group.ParentName.Valid {
			parentName := node.Group.ParentName.String
			if parent, ok := groupMap[parentName]; ok {
				parent.Children = append(parent.Children, node)
			}
		} else {
			root.Children = append(root.Children, node)
		}
	}

	// Build an ordering map from the sorted roots and sort the tree
	orderMap := buildOrderingMap(sortedRoots)
	sortChartTree(&root, orderMap)

	// Order the accounts within the tree
	err = sortAccountsInTree(&root)
	if err != nil {
		return nil, err
	}

	return &root, nil
}

func populateGroupMap(n *accountGroupNode, groupMap map[string]*ChartOfAccountsNode) {
	groupMap[n.group.Name] = &ChartOfAccountsNode{
		Group:    n.group,
		Children: []*ChartOfAccountsNode{},
		Accounts: []*Account{},
	}

	for _, child := range n.children {
		populateGroupMap(child, groupMap)
	}
}

func buildOrderingMap(sortedRoots []*accountGroupNode) map[string]int {
	orderMap := make(map[string]int)
	index := 0

	var traverse func(node *accountGroupNode)
	traverse = func(node *accountGroupNode) {
		orderMap[node.group.Name] = index
		index++
		for _, child := range node.children {
			traverse(child)
		}
	}

	for _, node := range sortedRoots {
		traverse(node)
	}

	return orderMap
}

func sortChartTree(root *ChartOfAccountsNode, orderMap map[string]int) {
	sort.SliceStable(root.Children, func(i, j int) bool {
		return orderMap[root.Children[i].Group.Name] < orderMap[root.Children[j].Group.Name]
	})

	for _, child := range root.Children {
		sortChartTree(child, orderMap)
	}
}

func sortAccountsInTree(node *ChartOfAccountsNode) error {
	idFn := func(a *Account) string { return a.Name }
	afterFn := func(a *Account) (string, bool) {
		if a.DisplayAfter.Valid {
			return a.DisplayAfter.String, true
		}
		return "", false
	}

	sorted, err := ordering.TopoSort[*Account](node.Accounts, idFn, afterFn)
	if err != nil {
		return err
	}

	node.Accounts = sorted

	for _, child := range node.Children {
		if err := sortAccountsInTree(child); err != nil {
			return err
		}
	}

	return nil
}
