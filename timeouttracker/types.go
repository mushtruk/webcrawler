package timeouttracker

import "sync"

type TimeoutTracker struct {
	mu   sync.Mutex
	urls map[string]struct{}
}

func NewTimeoutTracker() *TimeoutTracker {
	return &TimeoutTracker{
		urls: make(map[string]struct{}),
	}
}

func (t *TimeoutTracker) Add(url string) {
	t.mu.Lock()
	t.urls[url] = struct{}{}
	t.mu.Unlock()
}

func (t *TimeoutTracker) HasTimeout(url string) bool {
	t.mu.Lock()
	_, exists := t.urls[url]
	t.mu.Unlock()
	return exists
}

func (t *TimeoutTracker) Reset() {
	t.mu.Lock()
	t.urls = make(map[string]struct{})
	t.mu.Unlock()
}
