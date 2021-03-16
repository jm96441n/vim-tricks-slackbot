package slack

import (
	"fmt"
	"log"
	"regexp"
	"vim-tricks-slackbot/rss"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

type Bot struct{}

func NewBot() Bot {
	return Bot{}
}

func (b Bot) SendMessage(item *rss.Item) error {
	fmt.Println(item.Description)
	itemDesc, err := b.FormatMarkdown(item.Description)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("FORMATTED")
	fmt.Println(itemDesc)
	return nil
}

func (b Bot) FormatMarkdown(text string) (string, error) {
	markdown, err := b.toMarkdown(text)
	if err != nil {
		return "", err
	}
	fmt.Println("MARKDOWN")
	fmt.Println(markdown)
	markdown, err = b.formatLinks(markdown)

	if err != nil {
		return "", err
	}
	return markdown, nil
}

func (b Bot) toMarkdown(text string) (string, error) {

	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(text)
	if err != nil {
		return "", err
	}
	return markdown, nil
}

var pattern = regexp.MustCompile(`\[(?P<title>([^\]])*)\]\((?P<url>https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*))\)`)
var replacer = regexp.MustCompile(`\[[^\]]*\]\(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)\)`)
var template = []byte("<$url|$title>")

func (b Bot) formatLinks(text string) (string, error) {
	results := make([]string, 0)
	for _, submatches := range pattern.FindAllSubmatchIndex([]byte(text), -1) {
		// Apply the captured submatches to the template and append the output
		// to the result.
		res := []byte{}
		results = append(results, string(pattern.Expand(res, template, []byte(text), submatches)))
	}
	i := 0
	replacerFn := func(_ string) string {
		res := results[i]
		i++
		return res
	}

	return replacer.ReplaceAllStringFunc(text, replacerFn), nil
}
