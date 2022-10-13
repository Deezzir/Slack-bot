package bots

import (
	"context"
)

type Bot interface {
	init()
	Start(ctx context.Context)

	setCommands(commands []Command)
}

type singleton struct {
	Instance Bot
}

var (
	slackBot *singleton
)
