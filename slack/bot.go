package slack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"vim-tricks-slackbot/rss"
)

type slackFormatter interface {
	FormatHTMLToMarkdown(string) (string, error)
}

type httpInteractor interface {
	Do(req *http.Request) (*http.Response, error)
}

// Bot handles interaction with slack
type Bot struct {
	formatter    slackFormatter
	client       httpInteractor
	vimChannelId string
}

// NewBot is the constructor for a bot
func NewBot(formatter slackFormatter, client httpInteractor) Bot {
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
	jsonMsg, err := json.Marshal(message.Blocks)
	if err != nil {
		return err
	}
	values := url.Values{"blocks": {string(jsonMsg[:])}}
	url := fmt.Sprintf("https://slack.com/api/chat.postMessage?channel=%s&%s&text=TrickForToday&unfurl_links=false&unfurl_media=false", os.Getenv("VIM_CHANNEL_ID"), values.Encode())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SLACK_TOKEN")))
	resp, err := b.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed Request with: %s\n", resp)
		return fmt.Errorf("Failed Request with: %d", resp.StatusCode)
	}

	return nil
}

var section = "section"
var divider = "divider"
var mkdwn = "mrkdwn"

func (b Bot) buildMessage(item *rss.Item) (Message, error) {
	m := Message{Channel: os.Getenv("VIM_CHANNEL_ID")}
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
	Channel string
	Blocks  []Section `json:"blocks"`
}

type Section struct {
	Type string      `json:"type"`
	Text *SubSection `json:"text,omitempty"`
}

type SubSection struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}
