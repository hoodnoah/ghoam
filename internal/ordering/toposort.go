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
	idToIndexMap, err := makeIDToIndexMap(items, idFn)
	if err != nil {
		return nil, err
	}

	// make inDegrees, adjacency lists
	inDegreeAdjacency, err := makeInDegreesAdjacencyList(items, idFn, afterFn, idToIndexMap)
	if err != nil {
		return nil, err
	}
	adjacencyList := inDegreeAdjacency.adjacencyList
	inDegrees := inDegreeAdjacency.inDegrees

	// 3. Kahn's algorithm
	kahnTree, err := makeKahnTree(
		items,
		adjacencyList,
		inDegrees,
	)
	if err != nil {
		return nil, err
	}

	return kahnTree, nil
}

// maps items' identifiers to their index within the list of items provided.
// errors in the event of a duplicate entry
func makeIDToIndexMap[T any](items []T, idFn func(T) string) (map[string]int, error) {
	idMap := make(map[string]int, len(items))

	for index, item := range items {
		id := idFn(item)

		// check for duplicates by polling the map
		if _, dup := idMap[id]; dup {
			return nil, fmt.Errorf("duplicate ID detected: %s", id)
		}

		idMap[id] = index
	}

	return idMap, nil
}

// given a list of items, produce their adjacency list and inDegrees for Kahn's algorithm
func makeInDegreesAdjacencyList[T any](
	items []T,
	idFn func(T) string,
	afterFn func(T) (string, bool),
	idToIndexMap map[string]int,
) (
	*struct {
		adjacencyList [][]int
		inDegrees     []int
	}, error) {

	n := len(items)
	adjacencyList := make([][]int, n)
	inDegrees := make([]int, n)

	for index, item := range items {
		if afterID, ok := afterFn(item); ok {
			parentIndex, exists := idToIndexMap[afterID]
			if !exists {
				return nil, fmt.Errorf("toposort: unknown reference %q (for %q", afterID, idFn(item))
			}
			adjacencyList[parentIndex] = append(adjacencyList[parentIndex], index) // afterID -> current
			inDegrees[index]++
		}
	}

	return &struct {
		adjacencyList [][]int
		inDegrees     []int
	}{
		adjacencyList: adjacencyList,
		inDegrees:     inDegrees,
	}, nil
}

// kahn's algorithm
func makeKahnTree[T any](
	items []T,
	adjacencyList [][]int,
	inDegrees []int,
) ([]T, error) {
	var kahnTree []T

	queue := list.New() // create a list of nodes, all of which have in-degrees of 0
	for index := range len(items) {
		if inDegrees[index] == 0 {
			queue.PushBack(index)
		}
	}

	for queue.Len() > 0 {
		element := queue.Front()
		queue.Remove(element)
		value := element.Value.(int)

		kahnTree = append(kahnTree, items[value])

		for _, w := range adjacencyList[value] {
			inDegrees[w]--
			if inDegrees[w] == 0 {
				queue.PushBack(w)
			}
		}
	}

	if len(kahnTree) != len(items) {
		return nil, fmt.Errorf("toposort: cycle detected")
	}

	return kahnTree, nil
}
