package skiplist

import (
	"math/rand/v2"
)

type Key interface {
	int | float32 | float64
}

type Value interface {
	int | float32 | float64
}

type Node[K Key, V Value] struct {
	NodeKey   K
	NodeValue V

	NodePointers []*Node[K, V]
}

const max_level = 32

type Skiplist[K Key, V Value] struct {
	start *Node[K, V]
	end   *Node[K, V]
}

func NewSkipList[K Key, V Value]() *Skiplist[K, V] {

	start := &Node[K, V]{
		NodePointers: make([]*Node[K, V], 1),
	}

	end := &Node[K, V]{
		NodePointers: make([]*Node[K, V], 0),
	}

	start.NodePointers[0] = end

	return &Skiplist[K, V]{
		start: start,
		end:   end,
	}

}

func randomLevel() int {
	/* generates a random integer which is used to define the height of a new node, based on a coin flip simulation */

	level := 1
	for rand.Float32() <= 0.5 && level <= max_level {
		level++
	}
	return level
}

func (s *Skiplist[K, V]) Search(targetKey K) *Node[K, V] {
	/* search function: locates a key if it exists in the skiplist */

	currentPointer := s.start
	var resultPointer *Node[K, V] = nil

	for level := len(s.start.NodePointers) - 1; level >= 0; level-- {
		nextPointer := currentPointer.NodePointers[level]

		for nextPointer != s.end && nextPointer.NodeKey <= targetKey {
			currentPointer = nextPointer
			nextPointer = currentPointer.NodePointers[level]
		}

		if currentPointer != s.start && currentPointer.NodeKey == targetKey {
			resultPointer = currentPointer
			break
		}
	}

	return resultPointer
}

func (s *Skiplist[K, V]) Insert(Key K, Value V) *Node[K, V] {
	/* insert function: inserts a new node and returns pointer to it */

	if s.Search(Key) != nil {
		return nil
	}

	nodeHeight := randomLevel()
	var updateList []*Node[K, V]

	currentPointer := s.start

	for level := len(s.start.NodePointers) - 1; level >= 0; level-- {
		nextPointer := currentPointer.NodePointers[level]

		for nextPointer != s.end && nextPointer.NodeKey < Key {
			currentPointer = nextPointer
			nextPointer = currentPointer.NodePointers[level]
		}
		// next pointer's key is bigger
		updateList = append(updateList, currentPointer)
	}

	// construct the new node
	newNode := &Node[K, V]{
		NodeKey:      Key,
		NodeValue:    Value,
		NodePointers: make([]*Node[K, V], nodeHeight),
	}

	for ind, pointer := range updateList {
		level := len(updateList) - 1 - ind
		if nodeHeight > level {
			newNode.NodePointers[level] = pointer.NodePointers[level]
			pointer.NodePointers[level] = newNode
		}
	}

	for ind, _ := range newNode.NodePointers {
		if ind >= len(updateList) {
			s.start.NodePointers = append(s.start.NodePointers, newNode)
			newNode.NodePointers[ind] = s.end
		}
	}

	return newNode
}

func (s *Skiplist[K, V]) Delete(key K) bool {
	/* delete function: deletes a node from the skiplist and returns true if its succesfully deleted, false otherwise */

	var updateList []*Node[K, V]
	currentPointer := s.start

	for level := len(s.start.NodePointers) - 1; level >= 0; level-- {
		nextPointer := currentPointer.NodePointers[level]

		for nextPointer != s.end && nextPointer.NodeKey < key {
			currentPointer = nextPointer
			nextPointer = currentPointer.NodePointers[level]
		}
		updateList = append(updateList, currentPointer)
	}

	targetNode := currentPointer.NodePointers[0]

	if targetNode == s.end || targetNode.NodeKey != key {
		return false
	}

	for ind, pointer := range updateList {
		level := len(updateList) - 1 - ind

		if pointer.NodePointers[level] == targetNode {
			pointer.NodePointers[level] = targetNode.NodePointers[level]
		}
	}

	for len(s.start.NodePointers) > 1 && s.start.NodePointers[len(s.start.NodePointers)-1] == s.end {
		s.start.NodePointers = s.start.NodePointers[:len(s.start.NodePointers)-1]
	}

	return true
}

func (s *Skiplist[K, V]) Update(key K, newValue V) bool {
	/* update function: updates the value for the key if it exists, returns true on succesful operaion, false otherwise */
	node := s.Search(key)
	if node != nil {
		node.NodeValue = newValue
		return true
	}
	return false
}
