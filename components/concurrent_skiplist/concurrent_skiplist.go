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
	start          *SkipListNode // start node
	end            *SkipListNode // end node
	currentVersion atomic.Int64  // atomic counter for versioning
}

func NewSkipList() *SkipList {

	start := &SkipListNode{
		pointerArray: make([]atomic.Pointer[SkipListNode], 0),
	}

	end := &SkipListNode{}

	start.pointerArray[0].Store(end)

	var s1 atomic.Int64
	s1.Store(1)

	return &SkipList{
		start:          start,
		end:            end,
		currentVersion: s1,
	}
}

func randomLevel() int {
	/* generates a random integer which is used to define the height of a new node, based on a coin flip simulation */

	level := 1
	for rand.Float32() <= 0.5 && level <= 32 {
		level++
	}
	return level
}

func (s *SkipList) Seek(key int64, readTransactionNumber int64) *SkipListNode {
	// receives a key along with a read transaction number assigned to this request
	// returns a pointer to a node with the same key and latest possible versionNumber, else nil

	currentPointer := s.start // reads start, since start's pointer never changes we need not to use atomic here

	compositeTargetKey := CompositeKey{
		key:           key,
		versionNumber: readTransactionNumber,
	}

	for level := len(s.start.pointerArray) - 1; level >= 0; level-- {
		nextPointer := currentPointer.pointerArray[level].Load() // atomically load the forward pointer of current node for this level

		for nextPointer != s.end && (nextPointer.nodeKey.LessThan(&compositeTargetKey)) {
			currentPointer = nextPointer
			nextPointer = currentPointer.pointerArray[level].Load()
		}
	}

	// since we used lessThan we found the largest node which is less than compositeTargetKey, so the possible candidate is it's first child on level 0
	candidate := currentPointer.pointerArray[0].Load()

	if candidate != s.end && candidate.nodeKey.key == key {
		if !candidate.tombstone {
			return candidate
		}
	}

	return nil
}

func (S *SkipList) Insert(key int64, value int64, writeTransactionNumber int64, delete bool) *SkipListNode {
	// receives a key, value and a write transaction number assigned to this insert request
	// receives a delete boolean, to decide weather to flag a node as tombstone or not
	// returns a pointer to the inserted node, else nil
	// this writeTransactionNumber is calculated as (SkipList.load(currentVersion) + 1)
}

func (S *SkipList) Update(key int64, value int64, updateTransactionNumber int64) *SkipListNode {
	// receives a key, value and an update transaction number assigned to this update request
	// returns a pointer to the inserted node, else nil
	// since we want this skip list to be append only, we nevere update inplace instead for updating
	// we just insert a higher version node for the same key with the different updated value
	// internally calls Insert function
	return S.Insert(key, value, updateTransactionNumber, false)
}

func (S *SkipList) Delete(key int64, deleteTransactionNumber int64) bool {
	// receives a key and a delete transaction number assigned to this delete request
	// returns boolen based on succes of the delete operation
	// since this skip list is append only, we never actually delete the entry but add a newer version of the same key with a tombstone indicating it's deleted status
	// internally calls Insert function

	if S.Insert(key, key, deleteTransactionNumber, true) != nil {
		return true
	} else {
		return false
	}
}
