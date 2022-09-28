package main

import (
	"context"

	"slack-bot/pkg/bots"
	"slack-bot/pkg/config"
	"slack-bot/pkg/utils"
)

func main() {
	bot := &bots.SlackBot{
		Name:     "Noxu-Bot",
		BotToken: config.SLACK_BOT_TOKEN,
		AppToken: config.SLACK_APP_TOKEN,
	}
	bot.Init()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func() {
		if r := recover(); r != nil {
			utils.ErrorLogger.Printf("Recovered from panic: %v", r)
		}
	}()

	bot.Start(ctx)
}
