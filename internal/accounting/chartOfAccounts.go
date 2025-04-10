package accounting

import "context"

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

func buildChartNode(n *accountGroupNode, groupMap map[string]*ChartOfAccountsNode) *ChartOfAccountsNode {
	node := groupMap[n.group.Name]
	for _, child := range n.children {
		node.Children = append(node.Children, buildChartNode(child, groupMap))
	}

	return node
}
