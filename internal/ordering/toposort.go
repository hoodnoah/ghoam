package ordering

import (
	"container/list"
	"fmt"
)

// TopoSort orders `items` so that every element appears **after**
// the element whose ID is returned by `afterFn` (if any).
//
// idFn(item) -> unique identifier
// afterFn(item) -> ID of the item that should appear before this item
//
// It returns an error if it encounters
//   - duplicate IDs
//   - references to non-existent IDs
//   - cycles
func TopoSort[T any](
	items []T,
	idFn func(T) string,
	afterFn func(T) (string, bool),
) ([]T, error) {
	n := len(items)
	id2idx := make(map[string]int, n) // detect duplicates and resolve references
	indeg := make([]int, n)           // Kahn's in-degree counter
	adj := make([][]int, n)           // adjacency list

	// 1. Map IDs -> index, detect duplicates
	for i, it := range items {
		id := idFn(it)
		if _, dup := id2idx[id]; dup {
			return nil, fmt.Errorf("toposort: duplicate id %q", id)
		}
		id2idx[id] = i
	}

	// 2. Build graph and in-degree table
	for i, it := range items {
		if afterID, ok := afterFn(it); ok {
			parentIdx, exists := id2idx[afterID]
			if !exists {
				return nil, fmt.Errorf("toposort: unknown reference %q (for %q)", afterID, idFn(it))
			}
			adj[parentIdx] = append(adj[parentIdx], i) // afterID -> current
			indeg[i]++
		}
	}

	// 3. Kahn's algorithm
	q := list.New() // queue of nodes with in-degree 0
	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			q.PushBack(i)
		}
	}

	var out []T
	for q.Len() > 0 {
		e := q.Front()
		q.Remove(e)
		v := e.Value.(int)

		out = append(out, items[v])

		for _, w := range adj[v] {
			indeg[w]--
			if indeg[w] == 0 {
				q.PushBack(w)
			}
		}
	}

	if len(out) != n {
		return nil, fmt.Errorf("toposort: cycle detected")
	}

	return out, nil
}
