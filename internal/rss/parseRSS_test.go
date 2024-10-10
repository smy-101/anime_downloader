package rss

import (
	"testing"
)

func TestParseRss(t *testing.T) {
	rssUrl := "https://mikan.yujiangqaq.com/RSS/Bangumi?bangumiId=3403&subgroupid=615"

	items, err := ParseRSS(rssUrl)

	// 检查是否有错误
	if err != nil {
		t.Fatalf("ParseRSS returned an error: %v", err)
	}

	// 检查是否返回了项目
	if len(items) == 0 {
		t.Error("No items were parsed from the RSS feed")
	}

	// 检查每个项目的基本结构
	for i, item := range items {
		if item.Title == "" {
			t.Errorf("Item %d has an empty title", i)
		}
		if item.Link == "" {
			t.Errorf("Item %d has an empty link", i)
		}
		if item.torrent == "" {
			t.Errorf("Item %d has an empty torrent", i)
		}
		// Description 可能为空，所以我们不检查它
	}

	//打印items
	for _, item := range items {
		t.Logf("Title: %s\n", item.Title)
		t.Logf("Link: %s\n", item.Link)
		t.Logf("Description: %s\n", item.Description)
		t.Logf("Torrent: %s\n", item.torrent)
	}

}

func TestParseRSSNetworkError(t *testing.T) {
	// 使用一个不存在的 URL 来模拟网络错误
	_, err := ParseRSS("http://nonexistent.example.com")

	if err == nil {
		t.Error("Expected an error for non-existent URL, but got nil")
	}
}
