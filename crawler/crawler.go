package crawler

import (
	"github.com/mushtruk/webcrawler/fetcher"
	"github.com/mushtruk/webcrawler/parser"
	"github.com/mushtruk/webcrawler/queue"
)

type Crawler struct {
	Queue    *queue.Queue[*CrawlURL]
	Visited  map[string]bool
	MaxDepth int
}

func NewCrawler(q *queue.Queue[*CrawlURL], maxDepth int) *Crawler {
	return &Crawler{
		Queue:    q,
		Visited:  make(map[string]bool),
		MaxDepth: maxDepth,
	}
}

func (c *Crawler) Start() {
	for !c.Queue.IsEmpty() {
		crawlItem, ok := c.Queue.Next()
		if !ok {
			continue
		}

		url := crawlItem.ParsedURL.String()

		// Check if the URL has already been visited or if the depth limit has been reached
		if c.IsVisited(url) || crawlItem.Depth > c.MaxDepth {
			continue
		}

		c.MarkVisited(url)

		// Fetch and parse the URL's content
		content, err := fetcher.FetchUrl(url)
		if err != nil {
			// Handle error (e.g., log it)
			continue
		}

		newUrls, err := parser.ParseContent(content, url)
		if err != nil {
			// Handle error
			continue
		}

		// Add new URLs to the queue
		for _, newUrl := range newUrls {
			newCrawlURL, err := NewCrawlURL(newUrl, crawlItem.Depth+1)
			if err != nil {
				// Handle error
				continue
			}
			c.Queue.Add(newCrawlURL)
		}
	}
}

// MarkVisited marks a URL as visited.
func (c *Crawler) MarkVisited(url string) {
	c.Visited[url] = true
}

// IsVisited checks if a URL has been visited.
func (c *Crawler) IsVisited(url string) bool {
	_, visited := c.Visited[url]
	return visited
}
