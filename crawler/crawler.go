package crawler

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/mushtruk/webcrawler/parser"
	"github.com/mushtruk/webcrawler/treesitemap"
)

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

	c.Storage.Handle(c.RootNode)
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

	content, err := c.FetchUrl(url)

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

	if c.RootNode == nil {
		log.Println("RootNode is nil")
		return
	}

	node := &treesitemap.TreeNode{URL: url, Content: newUrls}

	c.RootNode.AddChild(node)
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

func (c *Crawler) FetchUrl(url string) (body string, err error) {

	if cachedBody, found := c.Cache.Get(url); found {
		return cachedBody, nil
	}

	client := http.Client{
		Timeout: time.Duration(c.TimeOut),
	}

	resp, err := client.Get(url)

	if err != nil {
		if os.IsTimeout(err) {
			c.TimeoutTracker.Add(url)
		}
		return "", err
	}
	defer resp.Body.Close()

	// Read and convert the body to string
	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	content := string(bodyBytes)

	c.Cache.Set(url, content)

	return content, nil
}
