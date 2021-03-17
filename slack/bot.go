package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"vim-tricks-slackbot/rss"
)

type slackFormatter interface {
	FormatHTMLToMarkdown(string) (string, error)
}

type httpPoster interface {
	Post(string, string, io.Reader) (*http.Response, error)
}

// Bot handles interaction with slack
type Bot struct {
	formatter slackFormatter
	client    httpPoster
}

// NewBot is the constructor for a bot
func NewBot(formatter slackFormatter, client httpPoster) Bot {
	return Bot{
		formatter: formatter,
		client:    client,
	}
}

// SendMessage formats the message and sends it so slack
func (b Bot) SendMessage(item *rss.Item) error {
	message, err := b.buildMessage(item)
	if err != nil {
		return err
	}
	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonMsg))
	resp, err := b.client.Post("localhost:8000", "application/json", bytes.NewBuffer(jsonMsg))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(b))
	}
	return nil
}

var section = "section"
var divider = "divider"
var mkdwn = "mrkdwn"

func (b Bot) buildMessage(item *rss.Item) (Message, error) {
	m := Message{}
	fmtTitle, err := b.formatter.FormatHTMLToMarkdown(item.Title)
	if err != nil {
		return m, err
	}
	fmtDesc, err := b.formatter.FormatHTMLToMarkdown(item.Description)
	if err != nil {
		return m, err
	}
	fmtContent, err := b.formatter.FormatHTMLToMarkdown(item.Content)
	if err != nil {
		return m, err
	}
	m.Blocks = []Section{
		Section{
			Type: section,
			Text: &SubSection{
				Type: mkdwn,
				Text: fmtTitle,
			},
		},
		Section{
			Type: divider,
		},
		Section{
			Type: section,
			Text: &SubSection{
				Type: mkdwn,
				Text: fmtDesc,
			},
		},
		Section{
			Type: section,
			Text: &SubSection{
				Type: mkdwn,
				Text: fmtContent,
			},
		},
	}

	return m, nil
}

type Message struct {
	Blocks []Section `json:"blocks"`
}

type Section struct {
	Type string      `json:"type"`
	Text *SubSection `json:"text,omitempty"`
}

type SubSection struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}
