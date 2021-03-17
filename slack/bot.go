package slack

import (
	"fmt"
	"log"
	"vim-tricks-slackbot/rss"
)

type slackFormatter interface {
	FormatHTMLToMarkdown(string) (string, error)
}

// Bot handles interaction with slack
type Bot struct {
	formatter slackFormatter
}

func NewBot(formatter slackFormatter) Bot {
	return Bot{formatter: formatter}
}

func (b Bot) SendMessage(item *rss.Item) error {
	itemDesc, err := b.formatter.FormatHTMLToMarkdown(item.Description)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println(itemDesc)
	return nil
}
