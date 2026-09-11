package skiplist

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

type Skiplist[K Key, V Value] struct {
	start *Node[K, V] // consider start to be the smallest key
	end   *Node[K, V] // end to be the larget key
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

func insert() {
}

func delete() {
}
func update() {
}
