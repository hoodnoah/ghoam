package accounting

import (
	"fmt"

	"github.com/hoodnoah/ghoam/internal/ordering"
)

// helper type for tree
type accountGroupNode struct {
	group    *AccountGroup
	children []*accountGroupNode
}

// hierarchicalSortAccountGroups performs a hierarchical sort on AccountGroups.
// It first builds a tree based on ParentName, then sorts each set of siblings using their DisplayAfter field.
// Finally, it performs a pre-order traversal to produce a flattened, ordered slice.
func hierarchicalSortAccountGroups(groups []*AccountGroup) ([]*AccountGroup, error) {
	// 1. Build a map of group names to nodes
	nodeMap := make(map[string]*accountGroupNode, len(groups))
	for _, group := range groups {
		// create a node for each group,
		// each node has no children to start
		nodeMap[group.Name] = &accountGroupNode{
			group:    group,
			children: nil,
		}
	}

	// 2. Build the tree structure: set up parent-child relationships
	var roots []*accountGroupNode
	for _, node := range nodeMap {
		// if the node has a ParentName,
		// try to find it.
		if node.group.ParentName.Valid {
			parent, exists := nodeMap[node.group.ParentName.String]
			// non-extant but specified parent node is an error condition
			if !exists {
				return nil, fmt.Errorf("unknown parent %q for group %q", node.group.ParentName.String, node.group.Name)
			}

			// add the current node to its parents children
			parent.children = append(parent.children, node)
		} else {
			// this node has no parent, ergo it's a root
			roots = append(roots, node)
		}
	}

	// 3. Recursively sort the children at each level using TopoSort
	var sortChildren func(n *accountGroupNode) error
	sortChildren = func(n *accountGroupNode) error {
		if len(n.children) > 0 {
			// Use TopoSort to order the children w/ the DisplayAfter field.
			// Here we need to sort []*accountGroupNode.
			sorted, err := ordering.TopoSort(n.children,
				func(n *accountGroupNode) string { return n.group.Name },
				func(n *accountGroupNode) (string, bool) {
					if n.group.DisplayAfter.Valid {
						return n.group.DisplayAfter.String, true
					}
					return "", false
				},
			)
			if err != nil {
				return err
			}
			n.children = sorted

			// Recursively sort children of these nodes
			for _, child := range n.children {
				if err := sortChildren(child); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Optionally, sort the roots themselves with TopoSort
	sortedRoots, err := ordering.TopoSort(roots,
		func(n *accountGroupNode) string { return n.group.Name },
		func(n *accountGroupNode) (string, bool) {
			if n.group.DisplayAfter.Valid {
				return n.group.DisplayAfter.String, true
			}
			return "", false
		},
	)
	if err != nil {
		return nil, err
	}

	// Apply sorting recursively on the sorted roots
	for _, root := range sortedRoots {
		if err := sortChildren(root); err != nil {
			return nil, err
		}
	}

	// 4. Flatten the tree via a pre-order traversal
	var result []*AccountGroup
	var traverse func(n *accountGroupNode)
	traverse = func(n *accountGroupNode) {
		result = append(result, n.group)
		for _, child := range n.children {
			traverse(child)
		}
	}

	// Traverse each sorted root
	for _, root := range sortedRoots {
		traverse(root)
	}

	return result, nil
}

func SortAccountGroupsInPlace(groups []*AccountGroup) error {
	sorted, err := hierarchicalSortAccountGroups(groups)
	if err != nil {
		return err
	}

	copy(groups, sorted)

	return nil
}

func sortAdjacentGroups(n *accountGroupNode) error {
	if len(n.children) > 0 {
		// Use TopoSort to order the children w/ the DisplayAfter field.
		// Here we need to sort []*accountGroupNode.
		sorted, err := ordering.TopoSort(n.children,
			func(n *accountGroupNode) string { return n.group.Name },
			func(n *accountGroupNode) (string, bool) {
				if n.group.DisplayAfter.Valid {
					return n.group.DisplayAfter.String, true
				}
				return "", false
			},
		)
		if err != nil {
			return err
		}
		n.children = sorted

		// Recursively sort children of these nodes
		for _, child := range n.children {
			if err := sortAdjacentGroups(child); err != nil {
				return err
			}
		}
	}
	return nil
}
