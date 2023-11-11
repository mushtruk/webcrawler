package parser

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func ParseContent(content string, baseURL string) ([]string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err // Handle error if base URL is invalid
	}

	dom, err := html.Parse(strings.NewReader(content))

	if err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	urls = parseDomLink(dom, urls, base)

	return urls, nil
}

func parseDomLink(n *html.Node, urls []string, base *url.URL) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				href, err := url.Parse(a.Val)
				if err != nil {
					continue // Handle or log the error as per your requirement
				}
				resolvedURL := base.ResolveReference(href).String()
				urls = append(urls, resolvedURL)
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		urls = parseDomLink(c, urls, base)
	}
	return urls
}
