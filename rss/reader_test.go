package rss_test

import (
	"encoding/xml"
	"errors"
	"net/http"
	"os"
	"testing"
	"vim-tricks-slackbot/rss"
)

type MockHTTPClient struct {
	errorOnGet      bool
	errorOnXMLParse bool
}

var httpErr = errors.New("Bad Req")

func (m MockHTTPClient) Get(url string) (*http.Response, error) {
	if m.errorOnGet {
		return nil, httpErr
	}

	xml, err := os.Open("xml_test.rss")
	if err != nil {
		return nil, err
	}
	return &http.Response{Body: xml}, nil
}

func TestGetLatestItem(t *testing.T) {
	tests := []struct {
		name          string
		reader        *rss.Reader
		expectedItem  rss.Item
		expectedError error
	}{
		{
			name:          "Returns the latest item",
			reader:        rss.NewReader("fake_url", MockHTTPClient{errorOnGet: false}),
			expectedItem:  latestItem(),
			expectedError: nil,
		},
		{
			name:          "Fail to get the RSS field",
			reader:        rss.NewReader("fake_url", MockHTTPClient{errorOnGet: true}),
			expectedItem:  rss.Item{},
			expectedError: rss.ErrHTTPGet{Err: httpErr},
		},
	}

	// TODO: Add test cases around failures to parse XML and when feed is empty
	for _, test := range tests {
		item, err := test.reader.GetLatestItem()
		if err != test.expectedError {
			t.Errorf("Expected: %s\n, Got: %s\n", test.expectedError, err)
		}

		if test.expectedError != nil {
			continue
		}

		derefItem := *item
		// Compare fields for expected item
		if derefItem.XMLName != test.expectedItem.XMLName {
			t.Errorf("Expected XMLName: %s,\n Got: %s", test.expectedItem.XMLName, derefItem.XMLName)
		}

		if derefItem.Title != test.expectedItem.Title {
			t.Errorf("Expected XMLName: %s,\n Got: %s", test.expectedItem.Title, derefItem.Title)
		}

		if derefItem.Link != test.expectedItem.Link {
			t.Errorf("Expected Link: %s,\n Got: %s", test.expectedItem.Link, derefItem.Link)
		}

		if derefItem.PubDate != test.expectedItem.PubDate {
			t.Errorf("Expected PubDate: %s,\n Got: %s", test.expectedItem.PubDate, derefItem.PubDate)
		}

		if derefItem.Description != test.expectedItem.Description {
			t.Errorf("Expected Description: \n%q\n Got: \n%q\n", test.expectedItem.Description, derefItem.Description)
		}
	}
}

func latestItem() rss.Item {
	return rss.Item{
		XMLName: xml.Name{Local: "item"},
		Title:   "Copying and pasting lines",
		Link:    "https://vimtricks.com/p/copying-and-pasting-lines/",
		PubDate: "Thu, 11 Mar 2021 12:23:11 +0000",
		Description: `<p>Use ex commands in Vim to copy lines around without moving.</p>
<p>The post <a rel="nofollow" href="https://vimtricks.com/p/copying-and-pasting-lines/" data-wpel-link="internal">Copying and pasting lines</a> appeared first on <a rel="nofollow" href="https://vimtricks.com" data-wpel-link="internal">VimTricks</a>.</p>
`,
	}
}
