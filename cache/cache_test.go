package cache_test

import (
	"sync"
	"testing"

	"github.com/mushtruk/webcrawler/cache"
)

func TestCache(t *testing.T) {
	cache := cache.NewCache()

	url := "https://test-cache-url.com"
	expectedContent := "Test content"

	cache.Set(url, expectedContent)

	cachedContent, found := cache.Get(url)

	if !found {
		t.Errorf("Expected to get %s content for %s key but got nil", expectedContent, url)
	}

	if expectedContent != cachedContent {
		t.Errorf("Expected to get %s content but got %s", expectedContent, cachedContent)
	}
}

func TestCacheConcurrency(t *testing.T) {
	cache := cache.NewCache()
	url := "http://example.com"
	content := "Example Content"

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Set(url, content) // Simultaneous writes
			_, _ = cache.Get(url)   // Simultaneous reads
		}()
	}

	wg.Wait()
	// Check if cache has the correct value after concurrent access
	if cachedContent, _ := cache.Get(url); cachedContent != content {
		t.Errorf("Cache content mismatch after concurrent access")
	}
}
