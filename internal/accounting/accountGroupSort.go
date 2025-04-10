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

// hierarchicalAccountGroupSort performs a hierarchical sort on AccountGroups.
// It first builds a tree based on ParentName, then sorts each set of siblings using their DisplayAfter field.
// Finally, it performs a pre-order traversal to produce a flattened, ordered slice.
func hierarchicalAccountGroupSort(groups []*AccountGroup) ([]*AccountGroup, error) {
	// Get sorted tree of accountGroupNodes
	sortedRoots, err := BuildAccountGroupTree(groups)
	if err != nil {
		return nil, err
	}

	// Apply sorting recursively on the sorted roots
	for _, root := range sortedRoots {
		if err := sortAdjacentGroups(root); err != nil {
			return nil, err
		}
	}

	// 4. Flatten the tree via a pre-order traversal
	var result []*AccountGroup

	// Traverse each sorted root
	for _, root := range sortedRoots {
		result = traverse(root, result)
	}

	return result, nil
}

// builds and returns a sorted tree of accountGroupNodes.
func BuildAccountGroupTree(groups []*AccountGroup) ([]*accountGroupNode, error) {
	nodeMap := makeNodeMap(groups)
	roots, err := makeTree(nodeMap)
	if err != nil {
		return nil, err
	}

	sortedRoots, err := ordering.TopoSort(
		roots,
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

	for _, root := range sortedRoots {
		if err := sortAdjacentGroups(root); err != nil {
			return nil, err
		}
	}

	return sortedRoots, nil
}

// sorts account groups in place (pseudo, it copies them internally)
func SortAccountGroupsInPlace(groups []*AccountGroup) error {
	sorted, err := hierarchicalAccountGroupSort(groups)
	if err != nil {
		return err
	}

	copy(groups, sorted)

	return nil
}

// Creates a map of accountGroups as nodes, with their Name as the key
func makeNodeMap(groups []*AccountGroup) map[string]*accountGroupNode {
	nodeMap := make(map[string]*accountGroupNode, len(groups))
	for _, group := range groups {
		nodeMap[group.Name] = &accountGroupNode{
			group:    group,
			children: []*accountGroupNode{},
		}
	}

	return nodeMap
}

// Populates a tree of nodes with parent-child relationships
func makeTree(nodeMap map[string]*accountGroupNode) ([]*accountGroupNode, error) {
	var roots []*accountGroupNode
	for _, node := range nodeMap {
		// if the node has a ParentName, try to find the associated node
		if node.group.ParentName.Valid {
			parent, exists := nodeMap[node.group.ParentName.String]
			if !exists {
				// non-extant but specified parent node is an error condition
				// it means a child references a parent which doesn't exist
				return nil, fmt.Errorf("unknown parent %q for group %q", node.group.ParentName.String, node.group.Name)
			}

			// Add the current node to its parent's children
			parent.children = append(parent.children, node)
		} else {
			// this node has no valid ParentName; it's a root node
			roots = append(roots, node)
		}
	}

	return roots, nil
}

// traverses a tree, flattening it into a slice
func traverse(n *accountGroupNode, output []*AccountGroup) []*AccountGroup {
	output = append(output, n.group)
	for _, child := range n.children {
		output = traverse(child, output)
	}

	return output
}

// sorts adjacent groups by their DisplayAfter field
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
