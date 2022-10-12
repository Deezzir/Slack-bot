package bots

import (
	"context"
	"slack-bot/pkg/utils"
)

type Bot interface {
	init()
	Start(ctx context.Context)

	setCommands(commands []Command)
}

var (
	slackBot *utils.Singleton
)
