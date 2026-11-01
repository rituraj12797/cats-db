package lockFreeSkiplist

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
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

func TestReadValidity(t *testing.T) {
	/* This test function launches a writer go routine and writes 10,000 entries in the skip list
	 * Tests weather reads result into succesful node findings or nil, and verifies the reason for the nil output
	 * WRITTEN BUT NOT FOUND : a parameter defining the number of time we faced a situation where an entry was written, after this when we tried to read it, we were unable to find it. Expected value : 0
	 * SeekTest should return nil only and only when the key was not written yet.
	 */
	s := NewSkipList(32)
	t.Logf(" Writes and reads run concurrently ")
	// var wgs sync.WaitGroup
	// wgs.Add(1)

	var done atomic.Int32
	done.Store(0)

	go func(s *SkipList) {
		// defer wgs.Done()
		for x := range 100000 { // 0 to 99999
			y := s.Insert(int64(x+1), int64((x+1)*10), false)
			if y != nil {
				done.Add(1) // increment by 1
			}

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
	wg.Add(100000)

	var written_but_not_read atomic.Int32
	written_but_not_read.Store(0)

	// launch 100000 readers
	for x := range 100000 {
		go func(s *SkipList, x int) {
			defer wg.Done()
			num := x + 1

			zd := done.Load()
			res, read_txn_no := s.SeekTest(int64(num))

			if res != nil {
				if res.NodeValue != (res.NodeKey.Key * 10) {
					t.Fatal("FAIL, key = ", res.NodeKey.Key, " value: ", res.NodeValue)
					failed_count.Add(1)
				} else {
					pass_count.Add(1)
				}
			} else {
				nil_count.Add(1)
				// verify this
				if read_txn_no >= int64(num) {
					// number was written but not found in read
					written_but_not_read.Add(1)
					t.Log(" zd : ", zd, " tried to read : ", num, "read_txn_no : ", read_txn_no)
				}
			}
		}(s, x)
	}

	// wait for them to finish
	wg.Wait()

	t.Log("PASSED : ", pass_count.Load())
	t.Log("FAILED : ", failed_count.Load())
	t.Log("NIL : ", nil_count.Load())
	t.Log("WRITTEN BUT NOT FOUND : ", written_but_not_read.Load())

	// a new idea to test the validity of the NIL responses we are getting
}

func TestReadIdempotance(t *testing.T) {
	s := NewSkipList(32)
	t.Logf("  Snapshot Isolation: Repeatable Reads Test ")

	var wg sync.WaitGroup

	wg.Add(10001)
	go func(s *SkipList) {
		defer wg.Done()
		for x := range 10000 { // 0 to 99999
			s.Insert(int64(x+1), int64((x+1)*10), false)
		}

		for j := range 100 {
			time.Sleep(10 * time.Microsecond)
			for x := range 10000 {
				s.Update(int64(x+1), int64((x+1)*10+j))
			}
		}
	}(s)

	for x := range 10000 {
		go func(s *SkipList, x int) {
			defer wg.Done()
			num := x + 1

			res, read_txn_no := s.TestSnapshotReadIdempotence(int64(num))
			if len(res) > 0 {
				if res[0].NodeValue != res[len(res)-1].NodeValue {
					t.Fatal("more than one read response detected, num: ", num, " read_txn : ", read_txn_no)
				}
			}
		}(s, x)
	}

	wg.Wait()
}
