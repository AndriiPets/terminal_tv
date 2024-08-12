package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

type SearchResult struct {
	Url   string
	Title string
}

func SearchYT(query string) ([]SearchResult, error) {

	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	// Creating timeout for 15 seconds
	ctx, cancel = context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	var videos []SearchResult

	var attributes []map[string]string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://www.youtube.com/results?search_query=decino`),

		chromedp.AttributesAll(`a#video-title`, &attributes, chromedp.ByQueryAll),
	)

	if err != nil {
		return videos, err
	}

	for _, attr := range attributes {
		url, ok := attr["href"]
		if !ok {
			return videos, fmt.Errorf("no url in attributes")
		}

		title, ok := attr["title"]
		if !ok {
			return videos, fmt.Errorf("no title in attributes")
		}
		vid := SearchResult{
			Url:   url,
			Title: title,
		}

		videos = append(videos, vid)
	}

	return videos, nil
}
