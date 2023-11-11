package queue

import (
	"strconv"
	"testing"
)

type IntKey int

func (ik IntKey) Key() string {
	return strconv.Itoa(int(ik)) // Convert int to string
}

// TestQueueAdd checks if items are correctly added to the queue.
func TestQueueAdd(t *testing.T) {
	q := NewQueue[IntKey]()
	q.Add(1)
	if q.IsEmpty() || q.Size() != 1 {
		t.Fatalf("Queue should have 1 item after adding, got %d", q.Size())
	}

	q.Add(2)
	if q.Size() != 2 {
		t.Fatalf("Queue should have 2 items after adding another, got %d", q.Size())
	}
}

// TestQueueNext checks if Next correctly returns items in FIFO order.
func TestQueueNext(t *testing.T) {
	q := NewQueue[IntKey]()
	q.Add(1)
	q.Add(2)

	first, ok := q.Next()
	if !ok || first != 1 {
		t.Fatalf("Expected first Next to return 1, got %d", first)
	}

	second, ok := q.Next()
	if !ok || second != 2 {
		t.Fatalf("Expected second Next to return 2, got %d", second)
	}

	if !q.IsEmpty() {
		t.Fatal("Queue should be empty after removing all items")
	}
}

// TestQueueIsEmpty checks if IsEmpty correctly reports the queue's status.
func TestQueueIsEmpty(t *testing.T) {
	q := NewQueue[IntKey]()
	if !q.IsEmpty() {
		t.Fatal("New queue should be empty")
	}

	q.Add(1)
	if q.IsEmpty() {
		t.Fatal("Queue should not be empty after adding an item")
	}

	_, _ = q.Next()
	if !q.IsEmpty() {
		t.Fatal("Queue should be empty after removing the item")
	}
}

func TestQueueStress(t *testing.T) {
	q := NewQueue[IntKey]()
	const numItems = 10000

	for i := 0; i < numItems; i++ {
		q.Add(IntKey(i))
	}

	if q.Size() != numItems {
		t.Errorf("Queue should have %d items, got %d", numItems, q.Size())
	}

	for i := 0; i < numItems; i++ {
		item, ok := q.Next()
		if !ok || item != IntKey(i) {
			t.Fatalf("Expected %d, got %d", i, item)
		}
	}

	if !q.IsEmpty() {
		t.Error("Queue should be empty after removing all items")
	}
}
