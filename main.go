package main

import (
	"context"

	"slack-bot/pkg/bots"
	"slack-bot/pkg/utils"
)

func main() {
	bot := bots.GetSlackBot().Instance.(*bots.SlackBot)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func() {
		if r := recover(); r != nil {
			utils.ErrorLogger.Printf("Recovered from panic: %v", r)
		}
	}()

	bot.Start(ctx)
}
