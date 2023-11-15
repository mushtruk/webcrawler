package crawler

import (
	"net/url"
	"sync"
	"time"

	"github.com/mushtruk/webcrawler/cache"
	"github.com/mushtruk/webcrawler/queue"
	"github.com/mushtruk/webcrawler/storage"
	"github.com/mushtruk/webcrawler/timeouttracker"
	"github.com/mushtruk/webcrawler/treesitemap"
)

// CrawlItem defines the interface for items that can be crawled.
type CrawlItem interface {
	GetDepth() int
	SetDepth(depth int)
}

// CrawlURL represents a URL to be crawled.
type CrawlURL struct {
	RawURL    string
	ParsedURL *url.URL
	Depth     int
}

type Crawler struct {
	Queue          *queue.Queue[*CrawlURL]
	Visited        map[string]bool
	MaxDepth       int
	TimeOut        time.Duration
	TimeoutTracker timeouttracker.TimeoutTracker
	Cache          cache.Cache
	Storage        storage.StorageHandler
	RootNode       *treesitemap.TreeNode
	mutex          sync.Mutex
}

func NewCrawler(q *queue.Queue[*CrawlURL], storage storage.StorageHandler, maxDepth int, timeout time.Duration) *Crawler {
	return &Crawler{
		Queue:          q,
		Visited:        make(map[string]bool),
		Cache:          *cache.NewCache(),
		RootNode:       &treesitemap.TreeNode{},
		TimeoutTracker: *timeouttracker.NewTimeoutTracker(),
		MaxDepth:       maxDepth,
		Storage:        storage,
		TimeOut:        timeout,
	}
}

// NewCrawlURL creates a new CrawlURL with the given raw URL and depth.
func NewCrawlURL(rawURL string, depth int) (*CrawlURL, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}
	return &CrawlURL{
		RawURL:    rawURL,
		ParsedURL: parsedURL,
		Depth:     depth,
	}, nil
}

func (c *CrawlURL) Key() string {
	return c.RawURL
}

// GetDepth returns the depth of the URL.
func (c *CrawlURL) GetDepth() int {
	return c.Depth
}

// SetDepth sets the depth of the URL.
func (c *CrawlURL) SetDepth(depth int) {
	c.Depth = depth
}

func (c *Crawler) MarkVisited(url string) {
	c.Visited[url] = true
}

func (c *Crawler) IsVisited(url string) bool {
	_, visited := c.Visited[url]
	return visited
}
