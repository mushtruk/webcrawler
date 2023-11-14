package cache

import (
	"sync"
)

type Cache struct {
	data map[string]string
	mu   sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		data: map[string]string{},
	}
}

func (c *Cache) Get(url string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	content, found := c.data[url]

	return content, found
}

func (c *Cache) Set(url, content string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[url] = content
}
