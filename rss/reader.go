package rss

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Reader is the struct responsible for providing functionality to read a RSS feed
type Reader struct {
	feedURL string
	client  httpGetter
}

type httpGetter interface {
	Get(string) (*http.Response, error)
}

// NewReader is the contstructor for the Reader struct
func NewReader(feedURL string, client httpGetter) *Reader {
	return &Reader{
		feedURL: feedURL,
		client:  client,
	}
}

// GetLatestItem returns the latest item from the RSS field
func (r *Reader) GetLatestItem() (*Item, error) {
	feed, err := r.buildFeed()
	if err != nil {
		return nil, err
	}

	item, err := feed.GetLatestItem()
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *Reader) buildFeed() (*Feed, error) {
	rawFeed, err := r.fetchFeed()

	if err != nil {
		return nil, ErrHTTPGet{Err: err}
	}

	return r.parseRawFeed(rawFeed)
}

func (r *Reader) fetchFeed() ([]byte, error) {
	resp, err := r.client.Get(r.feedURL)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (r *Reader) parseRawFeed(rawFeed []byte) (*Feed, error) {
	feed := &Feed{}
	decoder := xml.NewDecoder(bytes.NewReader(rawFeed))
	err := decoder.Decode(&feed)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

// ErrHTTPGet is when there is a failure to get the RSS Feed
type ErrHTTPGet struct {
	Err error
}

func (e ErrHTTPGet) Error() string {
	return fmt.Sprintf("Failed to get RSS Feed: %s", e.Err)
}
