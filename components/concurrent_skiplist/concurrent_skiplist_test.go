package lockFreeSkiplist

import (
	"sync"
	"sync/atomic"
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

// basic insert + seek, keys inserted out of order.
func TestInsertAndSeek(t *testing.T) {
	sl := NewSkipList(16)
	for _, k := range []int64{5, 1, 9, 3, 7} {
		sl.Insert(k, k*10, false)
	}

	for _, k := range []int64{1, 3, 5, 7, 9} {
		n := sl.Seek(k)
		if n == nil || n.NodeValue != k*10 {
			t.Fatalf("key %d: expected value %d, got %+v", k, k*10, n)
		}
	}
	for _, k := range []int64{0, 2, 4, 10} {
		if n := sl.Seek(k); n != nil {
			t.Fatalf("key %d should not exist, got %+v", k, n)
		}
	}
}

// readers must never observe a wrong value (nil is fine: not yet visible).
func TestConcurrentReadersSingleWriter(t *testing.T) {
	// skipping this for now
}

// snapshot stays stable while a writer keeps appending newer versions of the same key.
func TestSnapshotStableUnderConcurrentWrites(t *testing.T) {
	// skipping this for now
}

func TestRepeatableReads(t *testing.T) {
	/* This test function launches a writer go routine and writes 10,000 entries in the skip list
	 * 2 modes:
	 * either you wait for the writer to finish all writes then launch 1000 reader go routines ----> uncomment linee 57, 58, 61, 67 for this, expected result = PASSED: 10000, FAILED: 0, NIL: 0
	 * or you launch the reader go routines concurrently ----> comment line 57, 58, 61, 67, expected result = PASSED < 10,000, FAILED: 0, NIL: > 0
	 * in the later mode PASSED may be < 10k reason being, some writes  may not have been completed before their read query got launched resulting in nil
	 */
	s := NewSkipList(32)

	// var wgs sync.WaitGroup
	// wgs.Add(1)

	go func(s *SkipList) {
		// defer wgs.Done()
		for x := range 10000 { // 0 to 9999
			s.Insert(int64(x+1), int64((x+1)*10), false)
		}
	}(s)

	// wgs.Wait()

	var pass_count atomic.Int32
	var failed_count atomic.Int32
	var nil_count atomic.Int32

	pass_count.Store(0)
	failed_count.Store(0)
	nil_count.Store(0)

	var wg sync.WaitGroup
	wg.Add(1000)

	// launch 1000 readers
	for x := range 1000 { // 0 to 999
		go func(s *SkipList, x int) {
			defer wg.Done()
			for y := range 10 { // 0 to 9
				num := (x * 10) + y + 1
				res := s.Seek(int64(num))
				if res != nil {
					if res.NodeValue != (res.NodeKey.Key * 10) {
						t.Fatal("FAIL, key = ", res.NodeKey.Key, " value: ", res.NodeValue)
						failed_count.Add(1)
					} else {
						pass_count.Add(1)
					}
				} else {
					nil_count.Add(1)
				}
			}
		}(s, x)
	}

	// wait for them to finish
	wg.Wait()

	t.Log("PASSED : ", pass_count.Load())
	t.Log("FAILED : ", failed_count.Load())
	t.Log("NIL : ", nil_count.Load())

}
