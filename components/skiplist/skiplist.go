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
	NodeKey   K // key of this node
	NodeValue V // value of this node

	NodePointers []*Node[K, V] // array of pointers to node
}

const max_level = 32 // max height of a node in the skip list

type Skiplist[K Key, V Value] struct {
	start *Node[K, V] // consider start to be the smallest key
	end   *Node[K, V] // end to be the larget key
}

func randomLevel() int {
	level := 1
	for rand.Float32() <= 0.5 && level <= max_level {
		level++
	}
	return level
}

func (s *Skiplist[K, V]) search(targetKey K) *Node[K, V] {

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

func (s *Skiplist[K, V]) insert(Key K, Value V) *Node[K, V] {

	/*
	* Find if this key already exists or not
	* Not found ======>
	*
	 */

	if s.search(Key) != nil {
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

func delete() {

}

func (s *Skiplist[K, V]) update(key K, newValue V) bool {
	node := s.search(key)
	if node != nil {
		node.NodeValue = newValue
		return true
	}
	return false
}
