package timeouttracker

import (
	"testing"
)

func TestNewTimeoutTracker(t *testing.T) {
	tracker := NewTimeoutTracker()
	if tracker == nil {
		t.Fatal("NewTimeoutTracker() = nil, want non-nil")
	}
	if len(tracker.urls) != 0 {
		t.Errorf("New tracker should have an empty map, got %v", tracker.urls)
	}
}

func TestAdd(t *testing.T) {
	tracker := NewTimeoutTracker()
	url := "http://example.com"
	tracker.Add(url)

	if _, exists := tracker.urls[url]; !exists {
		t.Errorf("Add() did not add the url to the map")
	}
}

func TestHasTimeout(t *testing.T) {
	tracker := NewTimeoutTracker()
	url := "http://example.com"
	tracker.Add(url)

	if !tracker.HasTimeout(url) {
		t.Errorf("HasTimeout() = false, want true for existing url")
	}

	if tracker.HasTimeout("http://notexist.com") {
		t.Errorf("HasTimeout() = true, want false for non-existing url")
	}
}

func TestReset(t *testing.T) {
	tracker := NewTimeoutTracker()
	url := "http://example.com"
	tracker.Add(url)
	tracker.Reset()

	if len(tracker.urls) != 0 {
		t.Errorf("Reset() did not clear the urls map")
	}
}
