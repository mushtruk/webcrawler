package crawler

import (
	"net/url"
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

// GetDepth returns the depth of the URL.
func (c *CrawlURL) GetDepth() int {
	return c.Depth
}

// SetDepth sets the depth of the URL.
func (c *CrawlURL) SetDepth(depth int) {
	c.Depth = depth
}

func (c *CrawlURL) Key() string {
	return c.RawURL
}
