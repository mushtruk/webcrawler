package queue_test

import (
	"testing"

	"github.com/mushtruk/webcrawler/queue"
)

type mockItem struct {
	value string
}

func (m mockItem) Key() string {
	return m.value
}

func TestQueueAddAndVisited(t *testing.T) {
	q := queue.NewQueue[mockItem]()
	item1 := mockItem{"item1"}

	q.Add(item1)
	if q.Size() != 1 {
		t.Errorf("Queue size should be 1, got %d", q.Size())
	}

	if !q.IsVisited(item1) {
		t.Error("Item1 should be marked as visited")
	}

	// Test adding a duplicate item
	q.Add(item1)
	if q.Size() != 1 {
		t.Errorf("Queue size should still be 1 after adding duplicate, got %d", q.Size())
	}
}
