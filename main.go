package main

import (
	"fmt"
	"time"

	"github.com/mushtruk/webcrawler/crawler"
	"github.com/mushtruk/webcrawler/queue"
)

func main() {
	q := queue.NewQueue[*crawler.CrawlURL]()
	u, err := crawler.NewCrawlURL("https://medium.com/", 1)

	if err != nil {
		fmt.Printf("Error parsing url, got %v", u)
	}

	q.Add(u)

	c := crawler.NewCrawler(q, 2, 5*time.Second)

	c.Start(1)

	fmt.Print(c.Visited)
}
