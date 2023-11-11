package fetcher

import (
	"io"
	"net/http"
)

func FetchUrl(url string) (body string, err error) {
	// Fetch the webpage
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and convert the body to string
	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	content := string(bodyBytes)

	return content, nil
}
