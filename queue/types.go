package queue

// KeyProvider defines an interface for items that can provide a key.
type KeyProvider interface {
	Key() string
}

// Queue represents a FIFO (first in, first out) queue.
type Queue[T KeyProvider] struct {
	items   []T
	visited map[string]bool
}

// NewQueue creates a new Queue instance.
func NewQueue[T KeyProvider]() *Queue[T] {
	return &Queue[T]{items: make([]T, 0), visited: make(map[string]bool)}
}
