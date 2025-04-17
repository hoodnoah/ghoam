package ordering

import (
	"reflect"
	"strings"
	"testing"
)

// testItem is a simple type for testing purposes.
type testItem struct {
	id    string // unique identifier for the item
	after string // if non-empty, the ID this item must come after
}

// TestTopoSort runs a series of sub-tests against the TopoSort implementation.
func TestTopoSort(t *testing.T) {
	tests := []struct {
		name  string
		items []testItem
		// expectedOrder is specified only for deterministic expected outcomes.
		// When nil, we simply check that all dependency constraints are satisfied.
		expectedOrder []string
		expectErr     bool
		errSubstring  string
	}{
		{
			name: "simple chain",
			items: []testItem{
				{id: "A", after: ""},
				{id: "B", after: "A"},
				{id: "C", after: "B"},
			},
			expectedOrder: []string{"A", "B", "C"},
			expectErr:     false,
		},
		{
			name: "no dependencies",
			items: []testItem{
				{id: "A", after: ""},
				{id: "B", after: ""},
				{id: "C", after: ""},
			},
			// With no dependencies the algorithm will process them in input order.
			expectedOrder: []string{"A", "B", "C"},
			expectErr:     false,
		},
		{
			name: "duplicate IDs",
			items: []testItem{
				{id: "A", after: ""},
				{id: "A", after: ""},
			},
			expectErr:    true,
			errSubstring: "duplicate ID detected",
		},
		{
			name: "unknown dependency",
			items: []testItem{
				{id: "A", after: ""},
				{id: "B", after: "X"}, // "X" does not exist among the items
			},
			expectErr:    true,
			errSubstring: "unknown reference",
		},
		{
			name: "cycle detection",
			items: []testItem{
				{id: "A", after: "B"},
				{id: "B", after: "A"},
			},
			expectErr:    true,
			errSubstring: "cycle detected",
		},
		{
			name: "complex dependency",
			items: []testItem{
				{id: "A", after: ""},
				{id: "B", after: "A"},
				{id: "C", after: "A"},
				{id: "D", after: "B"},
				{id: "E", after: "C"},
				{id: "F", after: ""}, // independent item
			},
			// In this case the exact ordering is not unique.
			// We will verify that for every item with a dependency, its parent appears earlier.
			expectedOrder: nil,
			expectErr:     false,
		},
		{
			name: "default account groups ordering",
			items: []testItem{
				{id: "Equity", after: "Liabilities"},
				{id: "Expense", after: "Revenue"},
				{id: "Revenue", after: "Equity"},
				{id: "Liabilities", after: "Assets"},
				{id: "Assets", after: ""},
			},
			// Expecting the unique valid order:
			expectedOrder: []string{"Assets", "Liabilities", "Equity", "Revenue", "Expense"},
			expectErr:     false,
		},
	}

	// For each test case, create id and after functions and invoke TopoSort.
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Define id function: return the ID field.
			idFn := func(item testItem) string {
				return item.id
			}
			// Define after function: if the item has a non-empty "after" field, return it.
			afterFn := func(item testItem) (string, bool) {
				if item.after != "" {
					return item.after, true
				}
				return "", false
			}

			// Call TopoSort.
			sorted, err := TopoSort(tc.items, idFn, afterFn)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected an error but got nil")
				} else if !strings.Contains(err.Error(), tc.errSubstring) {
					t.Errorf("expected error to contain %q, got %q", tc.errSubstring, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Extract the ordered IDs from the sorted result.
			sortedIDs := make([]string, len(sorted))
			indexMap := make(map[string]int, len(sorted))
			for i, item := range sorted {
				id := idFn(item)
				sortedIDs[i] = id
				indexMap[id] = i
			}

			// If an exact expected order is provided, verify against it.
			if tc.expectedOrder != nil {
				if !reflect.DeepEqual(sortedIDs, tc.expectedOrder) {
					t.Errorf("expected order %v, got %v", tc.expectedOrder, sortedIDs)
				}
			} else {
				// Otherwise, verify dependency constraints:
				// for each item that depends on another, ensure that
				// the dependency appears earlier in the sorted list.
				for _, item := range tc.items {
					if parentID, ok := afterFn(item); ok {
						childPos := indexMap[idFn(item)]
						parentPos, exists := indexMap[parentID]
						if !exists {
							t.Errorf("parent %s not found in sorted results", parentID)
						} else if parentPos >= childPos {
							t.Errorf("dependency violation: %q (at %d) should come before %q (at %d)",
								parentID, parentPos, idFn(item), childPos)
						}
					}
				}
			}

			// Extra check: in the account groups case, "Assets" must always be first
			if indexMap["Assets"] != 0 {
				t.Errorf("expected \"Assets\" to be first, but found at position %d", indexMap["Assets"])
			}
		})
	}
}
