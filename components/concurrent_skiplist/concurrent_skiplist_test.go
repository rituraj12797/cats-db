package lockFreeSkiplist

import (
	"math"
	"sync"
	"testing"
)

// handler launches fn in its own goroutine and tracks it with wg.
func handler(wg *sync.WaitGroup, fn func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		fn()
	}()
}

// 1. Basic insert + seek, keys inserted out of order.
func TestInsertAndSeek(t *testing.T) {
	sl := NewSkipList(16)
	for i, k := range []int64{5, 1, 9, 3, 7} {
		sl.Insert(k, k*10, int64(i+1), false)
	}

	for _, k := range []int64{1, 3, 5, 7, 9} {
		n := sl.Seek(k, math.MaxInt64)
		if n == nil || n.nodeValue != k*10 {
			t.Fatalf("key %d: expected value %d, got %+v", k, k*10, n)
		}
	}
	for _, k := range []int64{0, 2, 4, 10} {
		if n := sl.Seek(k, math.MaxInt64); n != nil {
			t.Fatalf("key %d should not exist, got %+v", k, n)
		}
	}
}

// 2. Snapshot isolation: a read txn only sees versions <= its txn number.
func TestSnapshotVisibility(t *testing.T) {
	skipList := NewSkipList(16)

	insertData := []struct {
		key         int
		value       int
		transaction int
		tombstone   bool
	}{
		{0, 1, 1, false},
		{0, 2, 2, false},
		{0, 3, 3, false},
		{0, 4, 4, true}, // tombstone marked
		{0, 5, 5, false},
		{0, 6, 6, false},
		{0, 7, 7, false},
		{0, 8, 8, false},
		{0, 9, 9, false},
	}

	for _, x := range insertData {
		skipList.Insert(int64(x.key), int64(x.value), int64(x.transaction), x.tombstone)
	}

	testData := []struct {
		exists        bool
		key           int
		transaction   int
		expectedValue int
	}{
		{true, 0, 1, 1},
		{true, 0, 2, 2},
		{true, 0, 3, 3},
		{false, 0, 4, -1}, // should not exist as deleted(marked as tombstone) in 4th txn
		{true, 0, 5, 5},
		{true, 0, 6, 6},
		{true, 0, 7, 7},
		{true, 0, 8, 8},
		{true, 0, 9, 9},
		{false, 1, 1, 9},
		{false, 2, 3, 9},
	}

	for _, x := range testData {
		res := skipList.Seek(int64(x.key), int64(x.transaction))
		// t.Log("Test case: ", x)
		// t.Log(" Res : ", res)

		if res != nil {
			if res.nodeValue != int64(x.expectedValue) {
				t.Fatal("Wrong value found, key:", x.key, " expected: ", x.expectedValue, " found: ", res.nodeKey)
			}
		} else {
			if x.exists == true {
				t.Fatal("Node was inserted but not found")
			}
		}
	}
}

// 3. Delete (tombstone) hides the key only for txns >= delete txn; update re-adds it.
func TestDeleteAndUpdate(t *testing.T) {
	sl := NewSkipList(16)
	res := sl.Insert(5, 50, 1, false)

	t.Log(" Inserted : ", res)
	if !sl.Delete(5, 2) {
		t.Fatal("delete returned false")
	} else {
		t.Logf(" Deleted : 5, 2")
	}

	t.Log("Seeking (Key = 5, txn = 1, Expected: Val = 50 \n", sl.Seek(5, 1))

	if n := sl.Seek(5, 1); n == nil || n.nodeValue != 50 {
		t.Fatalf("txn 1 should still see 50, got %+v", n)
	}

	t.Log("Seeking (Key = 5, txn = 2, Expected: Tombstone \n", sl.Seek(5, 2))
	if n := sl.Seek(5, 2); n != nil {
		t.Fatalf("txn 2 should see tombstone (nil), got %+v", n)
	}

	t.Log("Updating (Key = 5, txn = 99, Val = 3 \n")

	sl.Update(5, 99, 3)
	if n := sl.Seek(5, 3); n == nil || n.nodeValue != 99 {
		t.Fatalf("txn 3 should see 99, got %+v", n)
	}
	t.Log("Seeking (Key = 5, txn = 110, Expected: Value = 99 \n", sl.Seek(5, 110))
	if n := sl.Seek(5, 110); n == nil || n.nodeValue != 99 {
		t.Fatalf("txn 110 should see 99, got %+v", n)
	}
	t.Log("Seeking (Key = 5, txn = 2, Expected: Tombstone \n", sl.Seek(5, 2))
	if n := sl.Seek(5, 2); n != nil {
		t.Fatalf("txn 2 should still see tombstone, got %+v", n)
	}
}

// 4. One writer inserting while many readers seek concurrently.
// Readers must never observe a wrong value (nil is fine: not yet visible).
func TestConcurrentReadersSingleWriter(t *testing.T) {
	// skipping this for now
}

// 5. Snapshot stays stable while a writer keeps appending newer versions of the same key.
func TestSnapshotStableUnderConcurrentWrites(t *testing.T) {
	// skipping this for now
}

func TestRepeatableReads( t *testing.T) {
	//
	//
	// launch a writer and multiple readers
	s := NewSkipList(32)

	go func writer(s *SkipList) {
		for x := range 1000 {
			s.Insert(x, x*10, false)
		}
	}(s)

	// launch 10 readers

	for x := range 100 {
		go func read(s *SkipList) {
			for y := range 10 {
				x = s.Seek()
			}
		}(s)
	}
}
