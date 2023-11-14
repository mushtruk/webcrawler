package parser_test

import (
	"reflect"
	"testing"

	"github.com/mushtruk/webcrawler/parser"
)

func TestParseContent_Basic(t *testing.T) {
	htmlContent := `<html><body><a href="http://example.com">Example</a></body></html>`
	baseURL := "http://test.com"

	expectedUrls := []string{"http://example.com"}

	urls, err := parser.ParseContent(htmlContent, baseURL)
	if err != nil {
		t.Fatalf("ParseContent returned an error: %v", err)
	}

	if !reflect.DeepEqual(urls, expectedUrls) {
		t.Errorf("Expected URLs %v, got %v", expectedUrls, urls)
	}
}

func TestParseContent_RelativeAndAbsolute(t *testing.T) {
	htmlContent := `<html><body><a href="/relative">Relative</a><a href="http://example.com">Absolute</a></body></html>`
	baseURL := "http://test.com"

	expectedUrls := []string{"http://test.com/relative", "http://example.com"}

	urls, err := parser.ParseContent(htmlContent, baseURL)
	if err != nil {
		t.Fatalf("ParseContent returned an error: %v", err)
	}

	if !reflect.DeepEqual(urls, expectedUrls) {
		t.Errorf("Expected URLs %v, got %v", expectedUrls, urls)
	}
}

func TestParseContent_InvalidURLs(t *testing.T) {
	htmlContent := `<html><body><a href=":%invalid">Invalid</a></body></html>`
	baseURL := "http://test.com"

	urls, _ := parser.ParseContent(htmlContent, baseURL)

	if len(urls) != 0 {
		t.Error("Expected an error for invalid URL, but got none")
	}
}

func TestParseContent_NoLinks(t *testing.T) {
	htmlContent := `<html><body><p>No links here</p></body></html>`
	baseURL := "http://test.com"

	urls, err := parser.ParseContent(htmlContent, baseURL)
	if err != nil {
		t.Fatalf("ParseContent returned an error: %v", err)
	}

	if len(urls) != 0 {
		t.Errorf("Expected no URLs, but got %v", urls)
	}
}

func TestParseContent_MalformedHTML(t *testing.T) {
	htmlContent := `<html><body><a href="http://example.com">Example</a></body>`
	// Missing closing HTML tag
	baseURL := "http://test.com"

	_, err := parser.ParseContent(htmlContent, baseURL)
	if err != nil {
		t.Fatalf("Did not expect an error for malformed HTML, but got: %v", err)
	}
}
