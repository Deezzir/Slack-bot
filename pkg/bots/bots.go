package bots

import (
	"context"
	"slack-bot/pkg/utils"
	"strings"
	"time"

	"github.com/shomali11/slacker"
)

type Bot interface {
	Init()
	Start(ctx context.Context)

	setCommands(commands []Command)
}

type SlackBot struct {
	Name     string
	BotToken string
	AppToken string

	bot *slacker.Slacker
}

func (s *SlackBot) Init() {
	s.bot = slacker.NewClient(s.BotToken, s.AppToken)
	s.setCommands(SlackCommands...)

	s.setHandlers()
}

func (s *SlackBot) setHandlers() {
	s.bot.Init(func() {
		utils.InfoLogger.Println("Slack Bot is initializing")
	})

	s.bot.Err(func(err string) {
		utils.ErrorLogger.Println(err)
	})

	s.bot.DefaultCommand(func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
		response.Reply("I don't understand what you mean")
	})

	s.bot.CleanEventInput(func(in string) string {
		in = strings.ReplaceAll(in, "`", `"`)
		in = strings.ReplaceAll(in, "'", `"`)
		return in
	})
}

func (s *SlackBot) logEvents() {
	for event := range s.bot.CommandEvents() {
		utils.CommandLogger.Printf("BOT: %s ", s.Name)
		utils.CommandLogger.Printf("TIME: %s ", event.Timestamp.Format(time.UnixDate))
		utils.CommandLogger.Printf("COMMAND: %s ", event.Command)
		utils.CommandLogger.Printf("PARAMETERS: %s", event.Parameters)
		utils.CommandLogger.Printf("EVENT: %s\n\n", event.Event)
	}
}

func (s *SlackBot) setCommands(commands ...Command) {
	for _, command := range commands {
		name := command.(SlackCommand).Name
		definition := command.(SlackCommand).CommandDefinition
		s.bot.Command(name, definition)
	}
}

func (s *SlackBot) Start(ctx context.Context) {
	if s.bot != nil {
		go s.logEvents()

		if err := s.bot.Listen(ctx); err != nil {
			panic(err)
		}
	} else {
		panic("SlackBot was not initialized, run Init() first")
	}
}
