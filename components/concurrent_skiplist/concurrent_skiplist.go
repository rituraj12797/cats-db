package lockFreeSkiplist

import (
	"math/rand/v2"
	"sync/atomic"
)

// Node key is composite which consists of original key + transaction_id in which it was created
// Sorted in increasing order of original key, and for same actual keys sorted in decreasing order of version number

type CompositeKey struct {
	key           int64
	versionNumber int64
}

func (K *CompositeKey) LessThan(other *CompositeKey) bool {
	if K.key < other.key {
		return true
	}
	if K.key == other.key && K.versionNumber > other.versionNumber {
		return true
	}
	return false
}

func (K *CompositeKey) LessThanEqualTo(other *CompositeKey) bool {
	if K.key < other.key {
		return true
	}
	if K.key == other.key && K.versionNumber >= other.versionNumber {
		return true
	}
	return false
}

func (K *CompositeKey) MoreThan(other *CompositeKey) bool {
	if K.key > other.key {
		return true
	}
	if K.key == other.key && K.versionNumber < other.versionNumber {
		return true
	}
	return false
}

func (K *CompositeKey) MoreThanEqualTo(other *CompositeKey) bool {
	if K.key > other.key {
		return true
	}
	if K.key == other.key && K.versionNumber <= other.versionNumber {
		return true
	}
	return false
}

func (K *CompositeKey) EqualTo(other *CompositeKey) bool {
	return K.key == other.key && K.versionNumber == other.versionNumber
}

type SkipListNode struct {
	nodeKey      CompositeKey
	nodeValue    int64
	pointerArray []atomic.Pointer[SkipListNode] // array of atomic pointers
	tombstone    bool
}

type SkipList struct {
	start            *SkipListNode // start node
	end              *SkipListNode // end node
	currentVersion   atomic.Int64  // atomic counter for versioning
	currentMaxHeight atomic.Int32
	heightLimit      int
}

func NewSkipList(limit int) *SkipList {

	if limit < 5 { // minimum height of the skip list should be 5 atleast
		return nil
	}

	start := &SkipListNode{
		pointerArray: make([]atomic.Pointer[SkipListNode], limit), // define the start node with maximum height
	}

	end := &SkipListNode{}

	for i := range start.pointerArray {
		start.pointerArray[i].Store(end)
	}

	var s1 atomic.Int64
	s1.Store(1)

	var x atomic.Int32
	x.Store(1)

	return &SkipList{
		start:            start,
		end:              end,
		currentVersion:   s1,
		heightLimit:      limit,
		currentMaxHeight: x,
	}
}

func (s *SkipList) randomLevel() int {
	/* generates a random integer which is used to define the height of a new node, based on a coin flip simulation */

	level := 1
	for rand.Float32() <= 0.5 && level < s.heightLimit {
		level++
	}
	return level
}

func (s *SkipList) Seek(key int64) *SkipListNode {
	// returns a pointer to a node with the same key and latest possible versionNumber, else nil

	currentPointer := s.start // reads start, since start's pointer never changes we need not to use atomic here
	var candidate *SkipListNode
	compositeTargetKey := CompositeKey{
		key:           key,
		versionNumber: s.currentVersion.Load(),
	}

	for level := s.currentMaxHeight.Load() - 1; level >= 0; level-- {
		nextPointer := currentPointer.pointerArray[level].Load() // atomically load the forward pointer of current node for this level

		for nextPointer != s.end && (nextPointer.nodeKey.LessThan(&compositeTargetKey)) {
			currentPointer = nextPointer
			nextPointer = currentPointer.pointerArray[level].Load()
		}
		// at level 0 when we arrive at this point means that we found the largest node which is less than compositeTargetKey, so the possible candidate is it's first child on level 0
		candidate = nextPointer
	}

	if candidate != s.end && candidate.nodeKey.key == key {
		if !candidate.tombstone {
			return candidate
		}
	}

	return nil
}

func (s *SkipList) Insert(key int64, value int64, delete bool) *SkipListNode {
	// receives a key, value and a write transaction number assigned to this insert request
	// receives a delete boolean, to decide weather to flag a node as tombstone or not
	// returns a pointer to the inserted node, else nil
	// this writeTransactionNumber is calculated as (SkipList.load(currentVersion) + 1)

	// once complete return the pointer to the new node

	currentPointer := s.start

	compositeTargetKey := CompositeKey{
		key:           key,
		versionNumber: s.currentVersion.Load() + 1,
	}

	// step 1 : construct the backward and forward array of nodes
	backwardArray := make([]*SkipListNode, 0)                         // we will store the current pointer of all the nodes at each level
	for level := s.currentMaxHeight.Load() - 1; level >= 0; level-- { // start from the max height
		nextPointer := currentPointer.pointerArray[level].Load() // atomically load the forward pointer of current node for this level
		for nextPointer != s.end && (nextPointer.nodeKey.LessThan(&compositeTargetKey)) {
			currentPointer = nextPointer
			nextPointer = currentPointer.pointerArray[level].Load()
		}
		backwardArray = append(backwardArray, currentPointer) // now nextpointer is GreaterthanEqualTo(composite(key,writeTransactionNumber))
		// but since writeTransactionNumber is new write txn and does not exist in the skip list untill now, so next pointer is GreaterThan(composite(key,writeTransactionNumber))
	}

	// step 2 : create a node with CompositeKey(key, writeTransactionNumber) and value
	height := s.randomLevel()

	newNode := &SkipListNode{
		nodeKey:      compositeTargetKey,
		nodeValue:    value,
		tombstone:    delete,
		pointerArray: make([]atomic.Pointer[SkipListNode], height),
	}

	// using the backwardArray we can find all (levell, current_node, rejection_node)
	// step 3 : set the forward pointer for this newly created node
	for index, node := range backwardArray {
		level := (len(backwardArray) - 1) - index
		if level < height {
			newNode.pointerArray[level].Store(node.pointerArray[level].Load())
		}
	}

	// step 4 : atomically swap the pointers of backward nodes to this at their corrosponding levels
	for level := 0; level < len(backwardArray) && level < height; level++ { // height check was missing to make sure we only change upto levels which can accomodate upto heigh
		// ex height =3 can only acomodate level 0 1 and 2 and not not be used as a level 3 node
		pred := backwardArray[len(backwardArray)-1-level]
		pred.pointerArray[level].Store(newNode)
	}

	s.currentVersion.Add(1) // new node inserted, now we can update he version counter

	// once back pointers are also attached and insert is complete atomically update currentMaxHeight
	// why didn't we do it earlier, since the already existing nodes were not attacthed to this node,
	// increasing the currentMaxHeight before hand wouldl only mean we will styarting our search from
	// a bigger height and no intermediate node between start and end at that height
	if height > int(s.currentMaxHeight.Load()) {
		s.currentMaxHeight.Store(int32(height)) // update the max height
	}

	return newNode
}

func (s *SkipList) Update(key int64, value int64) *SkipListNode {
	// receives a key, value and an update transaction number assigned to this update request
	// returns a pointer to the inserted node, else nil
	// since we want this skip list to be append only, we nevere update inplace instead for updating
	// we just insert a higher version node for the same key with the different updated value
	// internally calls Insert function
	//
	// need to check existence of this node
	seek_res := s.Seek(key)
	if (seek_res) != nil && seek_res.tombstone == false {
		// if the node exists
		return s.Insert(key, value, false)
	}
	// else this node does not exist
	return nil
}

func (s *SkipList) Delete(key int64) (success bool, status int) {
	seek_res := s.Seek(key)
	if seek_res != nil && seek_res.tombstone == false {
		if s.Insert(key, key, true) != nil {
			return true, 1
		}
		return false, 1
	}
	return false, 0
}
