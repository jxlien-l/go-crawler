package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Result struct {
	Medias    []string
	Links     []string
	Crawlable []string
}

func main() {
	httpClient := http.Client{
		Transport: &http.Transport{},
	}
	// Start url
	url := "https://go.dev/dl/"
	res, err := httpClient.Get(url)
	if err != nil {
		fmt.Println("err: %w", err)
	}
	baseURL := res.Request.URL.Scheme + res.Request.URL.Host
	content, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("err2: %w", err)
	}

	r := regexp.MustCompile(`(?:href|src)\s*=\s*["'](([^"']+\.(png|jpe?g|gif|webp|svg|bmp))|([^"']+\.(mp4|mkv|avi|mov|webm|flv))|([^"']+\.(pdf|docx?|xlsx?|pptx?|txt|rtf|odt))|([^"']+))["']`)

	matches := r.FindAllStringSubmatch(string(content), -1)
	fmt.Println("ok")

	var results Result
	for _, url := range matches {
		if url[2] != "" {
			results.Medias = append(results.Medias, fullURL(baseURL, url[2]))
		}
		if url[8] != "" {
			urlH := fullURL(baseURL, url[8])
			if isCrawlable(urlH) {
				results.Crawlable = append(results.Links, urlH)
			} else {
				results.Links = append(results.Links, urlH)
			}
		}
	}

	fmt.Printf("%+v\n", results)
}

func fullURL(baseURL, url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	}

	return baseURL + "://" + url
}

func isCrawlable(url string) bool {
	r := regexp.MustCompile(`^(net|com|de|fr)$`)
	frags := strings.Split(url, ".")
	test := r.MatchString(frags[len(frags)-1])
	fmt.Println(test)
	return r.MatchString(frags[len(frags)-1])
}
