package crawler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mushtruk/webcrawler/crawler"
	"github.com/mushtruk/webcrawler/queue"
)

const content = `<html>
<head><title>Test Page</title></head>
<body>
    <a href="/link1">Link 1</a>
    <a href="/link2">Link 2</a>
</body>
</html>`

// startTestServer returns a new mock HTTP server that responds with a simple HTML page containing links.
func newTestServer() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, content)
	})
	return httptest.NewServer(handler)
}

// TestCrawlerStart tests the basic functionality of the Crawler's Start method.
func TestCrawlerStart(t *testing.T) {
	// Set up a mock server with a simple HTML response
	server := newTestServer()
	defer server.Close()

	// Create a new crawler with the mock server's URL
	startURL, _ := crawler.NewCrawlURL(server.URL, 3)
	q := queue.NewQueue[*crawler.CrawlURL]()
	q.Add(startURL)
	crawler := crawler.NewCrawler(q, 3)

	// Start the crawler
	crawler.Start()

	// Verify the crawler has processed the URL and added new URLs to the queue
	if !q.IsEmpty() {
		t.Error("Queue should be empty after crawling")
	}

	if !crawler.IsVisited(startURL.RawURL) {
		t.Error("Start URL should be marked as visited after crawling")
	}

}
