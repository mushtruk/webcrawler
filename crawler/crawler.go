package crawler

import (
	"sync"

	"github.com/mushtruk/webcrawler/fetcher"
	"github.com/mushtruk/webcrawler/parser"
	"github.com/mushtruk/webcrawler/queue"
)

type Crawler struct {
	Queue    *queue.Queue[*CrawlURL]
	Visited  map[string]bool
	MaxDepth int
	mutex    sync.Mutex
}

func NewCrawler(q *queue.Queue[*CrawlURL], maxDepth int) *Crawler {
	return &Crawler{
		Queue:    q,
		Visited:  make(map[string]bool),
		MaxDepth: maxDepth,
	}
}

func (c *Crawler) Start() {
	var wg sync.WaitGroup

	for {
		c.mutex.Lock()
		if c.Queue.IsEmpty() {
			c.mutex.Unlock()
			break
		}
		crawlItem, ok := c.Queue.Next()
		c.mutex.Unlock()

		if !ok {
			continue
		}

		wg.Add(1)
		go func(ci *CrawlURL) {
			defer wg.Done()

			url := ci.ParsedURL.String()

			// Synchronized check for visited URLs and depth
			c.mutex.Lock()
			if c.IsVisited(url) || ci.Depth > c.MaxDepth {
				c.mutex.Unlock()
				return
			}
			c.MarkVisited(url)
			c.mutex.Unlock()

			// Fetch and parse the URL's content
			content, err := fetcher.FetchUrl(url)
			if err != nil {
				// Handle error
				return
			}

			newUrls, err := parser.ParseContent(content, url)
			if err != nil {
				// Handle error
				return
			}

			// Synchronized add to queue
			c.mutex.Lock()
			for _, newUrl := range newUrls {
				newCrawlURL, err := NewCrawlURL(newUrl, ci.Depth+1)
				if err == nil && !c.IsVisited(newCrawlURL.RawURL) && newCrawlURL.Depth <= c.MaxDepth {
					c.Queue.Add(newCrawlURL)
				}
			}
			c.mutex.Unlock()
		}(crawlItem)
	}

	wg.Wait()
}
