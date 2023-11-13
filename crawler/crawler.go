package crawler

import (
	"log"
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

func (c *Crawler) Start(workerCount int) {
	var wg sync.WaitGroup
	wg.Add(workerCount) // Add for the number of workers

	// Start a fixed number of worker goroutines
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for {
				c.mutex.Lock()
				if c.Queue.IsEmpty() {
					c.mutex.Unlock()
					return // No more URLs, exit the goroutine
				}

				crawlItem, ok := c.Queue.Next()
				c.mutex.Unlock()

				if ok {
					c.processURL(crawlItem)
				}
			}
		}()
	}

	wg.Wait() // Wait for all workers to finish
}

func (c *Crawler) noMoreURLsToProcess() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Queue.IsEmpty()
}

func (c *Crawler) getNextURLToCrawl() (*CrawlURL, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Queue.Next()
}

func (c *Crawler) processURL(ci *CrawlURL) {

	url := ci.ParsedURL.String()

	if c.shouldSkipURL(url, ci.Depth) {
		return
	}

	content, err := fetcher.FetchUrl(url)
	if err != nil {
		log.Printf("Error fetching URL %s: %v", url, err)
		return
	}

	newUrls, err := parser.ParseContent(content, url)

	if err != nil {
		log.Printf("Error parsing content from %s: %v", url, err)
		return
	}

	c.addNewURLsToQueue(newUrls, ci.Depth)
}

func (c *Crawler) shouldSkipURL(url string, depth int) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.IsVisited(url) || depth > c.MaxDepth {
		return true
	}
	c.MarkVisited(url)
	return false
}

func (c *Crawler) addNewURLsToQueue(newUrls []string, currentDepth int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, newUrl := range newUrls {
		newCrawlURL, err := NewCrawlURL(newUrl, currentDepth+1)

		if err == nil && !c.IsVisited(newCrawlURL.RawURL) && newCrawlURL.Depth <= c.MaxDepth {
			c.Queue.Add(newCrawlURL)
		}
	}
}
