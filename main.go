package main

import (
	"context"
	"fmt"
	"net/http"
	"vim-tricks-slackbot/rss"
	"vim-tricks-slackbot/slack"

	"github.com/aws/aws-lambda-go/lambda"
)

// VimTricksFeedURL is the url for the vim tricks feed
const VimTricksFeedURL string = "https://vimtricks.com/feed/"

// MyEvent is the event triggered by the lambda, we don't actually care what's in it just that it occurred
type MyEvent struct{}

type handler func(context.Context, MyEvent) (string, error)

func setupHandler(rssReader *rss.Reader, bot slack.Bot) handler {
	return func(ctx context.Context, event MyEvent) (string, error) {
		item, err := rssReader.GetLatestItem()
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		bot.SendMessage(item)
		return "Completed", nil
	}
}

func main() {
	client := http.DefaultClient
	rssReader := rss.NewReader(VimTricksFeedURL, client)

	bot := slack.NewBot()
	handlerFunc := setupHandler(rssReader, bot)
	lambda.Start(handlerFunc)
}
