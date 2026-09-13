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

type SkipListNode struct {
	nodeKey      CompositeKey
	nodeValue    int64
	pointerArray []*SkipListNode
	tombstone    bool
}

type SkipList struct {
	start          *SkipListNode // start node
	end            *SkipListNode // end node
	currentVersion atomic.Int64  // atomic counter for versioning
}

func NewSkipList() *SkipList {

	start := &SkipListNode{
		pointerArray: make([]*SkipListNode, 0),
	}

	end := &SkipListNode{}

	start.pointerArray = append(start.pointerArray, end)

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

func (S *SkipList) Seek(key int64, readTransactionNumber atomic.Int64) *SkipListNode {
	// receives a key along with a read transaction number assigned to this request
	// returns a pointer to a node with the same key and latest possible versionNumber, else nil
}

func (S *SkipList) Insert(key int64, value int64, writeTransactionNumber atomic.Int64, delete bool) *SkipListNode {
	// receives a key, value and a write transaction number assigned to this insert request
	// receives a delete boolean, to decide weather to flag a node as tombstone or not
	// returns a pointer to the inserted node, else nil
	// this writeTransactionNumber is calculated as (SkipList.load(currentVersion) + 1)
}

func (S *SkipList) Update(key int64, value int64, updateTransactionNumber atomic.Int64) *SkipListNode {
	// receives a key, value and an update transaction number assigned to this update request
	// returns a pointer to the inserted node, else nil
	// since we want this skip list to be append only, we nevere update inplace instead for updating
	// we just insert a higher version node for the same key with the different updated value
	// internally calls Insert function
	return S.Insert(key, value, updateTransactionNumber, false)
}

func (S *SkipList) Delete(key int64, deleteTransactionNumber atomic.Int64) bool {
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
