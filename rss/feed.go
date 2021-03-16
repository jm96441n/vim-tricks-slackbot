package rss

import (
	"encoding/xml"
	"errors"
)

// ErrEmptyItemList is returned when the RSS feed is empty and there are no items to return
var ErrEmptyItemList = errors.New("The feed is empty")

// Feed represents the entire RSS feed
type Feed struct {
	XMLNAme xml.Name `xml:"rss"`
	Channel *Channel `xml:"channel"`
}

// Channel represents the RSS channel
type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Title   string   `xml:"title"`
	Items   []*Item  `xml:"item"`
}

// Item represents an entry in the RSS feed
type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	PubDate     string   `xml:"pubDate"`
	Description string   `xml:"description"`
}

// GetLatestItem returns the latest item from the feed if there are items in the list
func (f *Feed) GetLatestItem() (*Item, error) {
	if len(f.Channel.Items) < 1 {
		return nil, ErrEmptyItemList
	}
	return f.Channel.Items[0], nil
}
