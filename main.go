package main

import (
	"fmt"
	"time"

	"github.com/mushtruk/webcrawler/crawler"
	"github.com/mushtruk/webcrawler/queue"
	"github.com/mushtruk/webcrawler/storage"
)

func main() {
	q := queue.NewQueue[*crawler.CrawlURL]()
	u, err := crawler.NewCrawlURL("https://placeholder.com/", 2)

	if err != nil {
		fmt.Printf("Error parsing url, got %v", u)
	}

	q.Add(u)

	s := storage.NewJSONStorage("placeholder.json", true)

	c := crawler.NewCrawler(q, s, 2, 5*time.Second)

	c.Start(10)
}
