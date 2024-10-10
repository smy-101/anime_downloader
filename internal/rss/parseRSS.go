package rss

import (
	"github.com/mmcdole/gofeed"
)

type RSSItem struct {
	Title       string
	Link        string
	Description string
	torrent     string
}

func ParseRSS(rssUrl string) ([]RSSItem, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssUrl)

	if err != nil {
		return nil, err
	}

	var items []RSSItem
	for _, item := range feed.Items {
		rssItem := RSSItem{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			torrent:     item.Enclosures[0].URL,
		}
		items = append(items, rssItem)
	}

	return items, nil

}
