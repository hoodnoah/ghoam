package accounting

import (
	"context"

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

	// build sorted (flat) tree of account groups
	groupTree, err := BuildAccountGroupTree(groups)
	if err != nil {
		return nil, err
	}

	// convert account groups to a tree
	pseudoRoot := ChartOfAccountsNode{
		Children: convertGroupTreeToChartNodes(groupTree),
	}

	// Create a lookup map to attach accounts to nodes
	nodeLookup := make(map[string]*ChartOfAccountsNode)
	collectNodes(&pseudoRoot, nodeLookup)

	// Attach accounts to the corresponding group nodes.
	for _, account := range accounts {
		if node, ok := nodeLookup[account.ParentGroupName]; ok {
			node.Accounts = append(node.Accounts, account)
		}
	}

	// Recursively topologically sort accounts at each node
	if err := topoSortAccountsInTree(&pseudoRoot); err != nil {
		return nil, err
	}

	return &pseudoRoot, nil
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

// Given a list of AccountGroup nodes, form them into an ordered tree structure
func convertGroupTreeToChartNodes(groupNodes []*accountGroupNode) []*ChartOfAccountsNode {
	var chartNodes []*ChartOfAccountsNode

	// for all groupNodes, create a ChartOfAccounts node
	// ordering is preserved
	for _, groupNode := range groupNodes {
		node := &ChartOfAccountsNode{
			Group:    groupNode.group,
			Children: convertGroupTreeToChartNodes(groupNode.children),
			Accounts: []*Account{},
		}
		chartNodes = append(chartNodes, node)
	}

	return chartNodes
}

// Builds a map from group name -> ChartOfAccountsNode for quick lookups
func collectNodes(node *ChartOfAccountsNode, lookup map[string]*ChartOfAccountsNode) {
	// if the node has no Group, add the group to the map
	if node.Group != nil {
		lookup[node.Group.Name] = node
	}

	// recurse down for all children
	for _, child := range node.Children {
		collectNodes(child, lookup)
	}
}

// Apply a topological sort on the accounts at the node, and recursively on its children
func topoSortAccountsInTree(root *ChartOfAccountsNode) error {
	// Sort accounts using a topological sort
	sorted, err := ordering.TopoSort(root.Accounts, AccountIDFn, AccountAfterFn)
	if err != nil {
		return err
	}

	root.Accounts = sorted

	// Apply the same ordering process to child nodes
	for _, child := range root.Children {
		if err := topoSortAccountsInTree(child); err != nil {
			return err
		}
	}

	return nil
}
