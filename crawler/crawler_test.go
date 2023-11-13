package crawler_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/mushtruk/webcrawler/crawler"
	"github.com/mushtruk/webcrawler/queue"
)

func randInRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func newMockServer(maxDepth int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		depth := strings.Count(r.URL.Path, "/link")
		// Start HTML content
		fmt.Fprint(w, "<html><body>")

		if depth < maxDepth {
			for i := 1; i < randInRange(2, 15); i++ {
				// Construct the URL for the next depth
				hostWithPort := r.Host // This includes the host and port
				nextDepthURL := fmt.Sprintf("http://%s%s/link%v", hostWithPort, strings.TrimSuffix(r.URL.Path, "/"), i)

				fmt.Fprintf(w, `<a href="%s">Link at depth %d</a><br>`, nextDepthURL, depth+1)
			}
		}

		// End HTML content
		fmt.Fprint(w, "</body></html>")
	}))

	return server
}

// TestCrawlerStart tests the basic functionality of the Crawler's Start method.
func TestCrawlerStart(t *testing.T) {
	// Set up a mock server with a simple HTML response
	server := newMockServer(10)
	defer server.Close()

	// Create a new crawler with the mock server's URL
	startURL, _ := crawler.NewCrawlURL(server.URL, 3)
	q := queue.NewQueue[*crawler.CrawlURL]()
	q.Add(startURL)
	crawler := crawler.NewCrawler(q, 3)

	// Start the crawler
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		crawler.Start(10)
	}()
	wg.Wait()

	// Verify the crawler has processed the URL and added new URLs to the queue
	if !q.IsEmpty() {
		t.Error("Queue should be empty after crawling")
	}

	if !crawler.IsVisited(startURL.RawURL) {
		t.Error("Start URL should be marked as visited after crawling")
	}

}

func TestCrawlerDepthControl(t *testing.T) {
	maxDepth := 3

	server := newMockServer(maxDepth)
	defer server.Close()

	startURL, _ := crawler.NewCrawlURL(server.URL, 0)
	q := queue.NewQueue[*crawler.CrawlURL]()
	q.Add(startURL)
	c := crawler.NewCrawler(q, maxDepth)

	c.Start(10)

	for visitedURL := range c.Visited {
		parsedURL, err := url.Parse(visitedURL)
		if err != nil {
			t.Errorf("Failed to parse URL: %s, error: %v", visitedURL, err)
			continue
		}

		depth, _ := strconv.Atoi(parsedURL.Query().Get("depth"))
		if depth > maxDepth {
			t.Errorf("Crawler exceeded max depth for URL: %s", visitedURL)
		}
	}
}

func TestCrawlerHighVolumeURLs(t *testing.T) {
	maxDepth := 5

	server := newMockServer(maxDepth)
	defer server.Close()

	q := queue.NewQueue[*crawler.CrawlURL]()

	crawlUrl, _ := crawler.NewCrawlURL(server.URL, 0)
	q.Add(crawlUrl)

	c := crawler.NewCrawler(q, maxDepth)

	c.Start(1000)

	if !q.IsEmpty() {
		t.Errorf("Queue should be empty after crawling")
	}
}
