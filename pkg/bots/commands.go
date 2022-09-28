package bots

import "github.com/shomali11/slacker"

type Command interface {
}

type SlackCommand struct {
	Name              string
	CommandDefinition *slacker.CommandDefinition
}

var SlackCommands = []Command{
	slackGetAge,
	slackYouSuck,
	slackListFiles,
	slackValidateEmail,
	slackGetFile,

	slackStartTicTacToe,
}
