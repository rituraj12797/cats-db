package skiplist

import (
	"testing"
)

func newSkiplist[K Key, V Value]() *Skiplist[K, V] {
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

func TestInsert(t *testing.T) {
	t.Run("Insert single element", func(t *testing.T) {
		sl := newSkiplist[int, int]()
		node := sl.Insert(10, 100)

		if node == nil {
			t.Error("Expected non-nil node after Insert")
		}
		if node.NodeKey != 10 {
			t.Errorf("Expected key 10, got %d", node.NodeKey)
		}
		if node.NodeValue != 100 {
			t.Errorf("Expected value 100, got %d", node.NodeValue)
		}
	})

	t.Run("Insert multiple elements", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		keys := []int{5, 10, 3, 20, 15}
		values := []int{50, 100, 30, 200, 150}

		for i, key := range keys {
			node := sl.Insert(key, values[i])
			if node == nil {
				t.Errorf("Failed to Insert key %d", key)
			}
		}

		for i, key := range keys {
			node := sl.Search(key)
			if node == nil {
				t.Errorf("Key %d not found after Insert", key)
			} else if node.NodeValue != values[i] {
				t.Errorf("Expected value %d for key %d, got %d", values[i], key, node.NodeValue)
			}
		}
	})

	t.Run("Insert duplicate key", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		node1 := sl.Insert(10, 100)
		if node1 == nil {
			t.Fatal("First Insert should succeed")
		}

		node2 := sl.Insert(10, 200)
		if node2 != nil {
			t.Error("Expected nil when Inserting duplicate key")
		}

		found := sl.Search(10)
		if found.NodeValue != 100 {
			t.Errorf("Expected original value 100, got %d", found.NodeValue)
		}
	})

	t.Run("Insert with float keys", func(t *testing.T) {
		sl := newSkiplist[float64, int]()

		sl.Insert(1.5, 15)
		sl.Insert(2.7, 27)
		sl.Insert(0.5, 5)

		node := sl.Search(1.5)
		if node == nil || node.NodeValue != 15 {
			t.Error("Failed to Insert/Search float key")
		}
	})
}

func TestSearch(t *testing.T) {
	t.Run("Search in empty skiplist", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		node := sl.Search(10)
		if node != nil {
			t.Error("Expected nil when Searching empty skiplist")
		}
	})

	t.Run("Search existing element", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		sl.Insert(10, 100)
		sl.Insert(20, 200)
		sl.Insert(5, 50)

		node := sl.Search(10)
		if node == nil {
			t.Fatal("Expected to find key 10")
		}
		if node.NodeKey != 10 {
			t.Errorf("Expected key 10, got %d", node.NodeKey)
		}
		if node.NodeValue != 100 {
			t.Errorf("Expected value 100, got %d", node.NodeValue)
		}
	})

	t.Run("Search non-existing element", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		sl.Insert(10, 100)
		sl.Insert(20, 200)
		sl.Insert(30, 300)

		node := sl.Search(15)
		if node != nil {
			t.Error("Expected nil when Searching for non-existing key")
		}
	})

	t.Run("Search smallest element", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		sl.Insert(10, 100)
		sl.Insert(5, 50)
		sl.Insert(20, 200)

		node := sl.Search(5)
		if node == nil || node.NodeValue != 50 {
			t.Error("Failed to find smallest element")
		}
	})

	t.Run("Search largest element", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		sl.Insert(10, 100)
		sl.Insert(5, 50)
		sl.Insert(20, 200)

		node := sl.Search(20)
		if node == nil || node.NodeValue != 200 {
			t.Error("Failed to find largest element")
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("Delete from empty skiplist", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		result := sl.Delete(10)
		if result {
			t.Error("Expected false when deleting from empty skiplist")
		}
	})

	t.Run("Delete existing element", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		sl.Insert(10, 100)
		sl.Insert(20, 200)
		sl.Insert(5, 50)

		result := sl.Delete(10)
		if !result {
			t.Error("Expected true when deleting existing element")
		}

		node := sl.Search(10)
		if node != nil {
			t.Error("Element should be Deleted")
		}

		if sl.Search(20) == nil {
			t.Error("Other elements should remain after deletion")
		}
		if sl.Search(5) == nil {
			t.Error("Other elements should remain after deletion")
		}
	})

	t.Run("Delete non-existing element", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		sl.Insert(10, 100)
		sl.Insert(20, 200)

		result := sl.Delete(15)
		if result {
			t.Error("Expected false when deleting non-existing element")
		}
	})

	t.Run("Delete all elements", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		keys := []int{10, 20, 5, 15, 25}
		for _, key := range keys {
			sl.Insert(key, key*10)
		}

		for _, key := range keys {
			result := sl.Delete(key)
			if !result {
				t.Errorf("Failed to Delete key %d", key)
			}
		}

		for _, key := range keys {
			if sl.Search(key) != nil {
				t.Errorf("Key %d should be Deleted", key)
			}
		}
	})

	t.Run("Delete smallest element", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		sl.Insert(10, 100)
		sl.Insert(5, 50)
		sl.Insert(20, 200)

		result := sl.Delete(5)
		if !result {
			t.Error("Failed to Delete smallest element")
		}

		if sl.Search(5) != nil {
			t.Error("Smallest element should be Deleted")
		}
	})

	t.Run("Delete largest element", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		sl.Insert(10, 100)
		sl.Insert(5, 50)
		sl.Insert(20, 200)

		result := sl.Delete(20)
		if !result {
			t.Error("Failed to Delete largest element")
		}

		if sl.Search(20) != nil {
			t.Error("Largest element should be Deleted")
		}
	})

	t.Run("Delete and reInsert", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		sl.Insert(10, 100)
		sl.Delete(10)
		node := sl.Insert(10, 200)

		if node == nil {
			t.Error("Failed to reInsert after deletion")
		}

		found := sl.Search(10)
		if found == nil || found.NodeValue != 200 {
			t.Error("ReInserted value should be 200")
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Update existing element", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		sl.Insert(10, 100)

		result := sl.Update(10, 999)
		if !result {
			t.Error("Expected true when updating existing element")
		}

		node := sl.Search(10)
		if node == nil {
			t.Fatal("Element should still exist after Update")
		}
		if node.NodeValue != 999 {
			t.Errorf("Expected Updated value 999, got %d", node.NodeValue)
		}
	})

	t.Run("Update non-existing element", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		sl.Insert(10, 100)

		result := sl.Update(20, 200)
		if result {
			t.Error("Expected false when updating non-existing element")
		}
	})

	t.Run("Update in empty skiplist", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		result := sl.Update(10, 100)
		if result {
			t.Error("Expected false when updating in empty skiplist")
		}
	})

	t.Run("Update multiple times", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		sl.Insert(10, 100)

		sl.Update(10, 200)
		sl.Update(10, 300)
		sl.Update(10, 400)

		node := sl.Search(10)
		if node == nil || node.NodeValue != 400 {
			t.Error("Expected final value 400 after multiple Updates")
		}
	})

	t.Run("Update different elements", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		sl.Insert(10, 100)
		sl.Insert(20, 200)
		sl.Insert(30, 300)

		sl.Update(10, 111)
		sl.Update(20, 222)
		sl.Update(30, 333)

		if sl.Search(10).NodeValue != 111 {
			t.Error("Failed to Update first element")
		}
		if sl.Search(20).NodeValue != 222 {
			t.Error("Failed to Update second element")
		}
		if sl.Search(30).NodeValue != 333 {
			t.Error("Failed to Update third element")
		}
	})

	t.Run("Update with float values", func(t *testing.T) {
		sl := newSkiplist[int, float64]()

		sl.Insert(10, 1.5)

		result := sl.Update(10, 2.7)
		if !result {
			t.Error("Failed to Update with float value")
		}

		node := sl.Search(10)
		if node == nil || node.NodeValue != 2.7 {
			t.Error("Float value not Updated correctly")
		}
	})
}

func TestRandomLevel(t *testing.T) {
	t.Run("Random level bounds", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			level := randomLevel()
			if level < 1 || level > max_level {
				t.Errorf("Random level %d out of bounds [1, %d]", level, max_level)
			}
		}
	})
}

func TestSkiplistIntegration(t *testing.T) {
	t.Run("Complex operations sequence", func(t *testing.T) {
		sl := newSkiplist[int, int]()

		for i := 0; i < 20; i += 2 {
			sl.Insert(i, i*10)
		}

		for i := 0; i < 20; i += 2 {
			node := sl.Search(i)
			if node == nil || node.NodeValue != i*10 {
				t.Errorf("Failed to find or incorrect value for key %d", i)
			}
		}

		for i := 0; i < 10; i += 2 {
			sl.Update(i, i*100)
		}

		for i := 0; i < 10; i += 2 {
			node := sl.Search(i)
			if node == nil || node.NodeValue != i*100 {
				t.Errorf("Update failed for key %d", i)
			}
		}

		for i := 0; i < 20; i += 4 {
			sl.Delete(i)
		}

		for i := 0; i < 20; i += 2 {
			node := sl.Search(i)
			if i%4 == 0 {
				if node != nil {
					t.Errorf("Key %d should be Deleted", i)
				}
			} else {
				if node == nil {
					t.Errorf("Key %d should still exist", i)
				}
			}
		}
	})
}
