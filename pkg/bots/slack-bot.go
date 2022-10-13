package bots

import (
	"context"
	"slack-bot/pkg/config"
	"slack-bot/pkg/utils"
	"strings"
	"time"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

type SlackBot struct {
	Name     string
	BotToken string
	AppToken string

	bot *slacker.Slacker
}

func GetSlackBot() *singleton {
	if slackBot == nil {
		utils.BotLock.Lock()
		defer utils.BotLock.Unlock()

		if slackBot == nil {
			bot := &SlackBot{
				Name:     "Noxu-Bot",
				BotToken: config.SLACK_BOT_TOKEN,
				AppToken: config.SLACK_APP_TOKEN,
			}
			bot.init()
			slackBot = &singleton{Instance: bot}
		}
	}
	return slackBot
}

func (s *SlackBot) init() {
	s.bot = slacker.NewClient(s.BotToken, s.AppToken)
	s.setCommands(SlackCommands)

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

func (s *SlackBot) setCommands(commands []Command) {
	for _, command := range commands {
		name := command.(SlackCommand).Name
		definition := command.(SlackCommand).CommandDefinition
		s.bot.Command(name, definition)
	}
}

func (s *SlackBot) PostMessage(channel, pretext, text string) (string, error) {
	client := s.bot.Client()

	attachment := slack.Attachment{
		Pretext: pretext,
		Text:    text,
		Color:   "#174dbe",
	}

	_, timestamp, err := client.PostMessage(
		channel,
		//slack.MsgOptionText("New message from bot", false),
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		utils.ErrorLogger.Printf("Failed to post a message - %s\n", err)
		return "", err
	}
	utils.InfoLogger.Printf("Message successfully sent to channel (%s)\n", channel)
	return timestamp, nil
}

func (s *SlackBot) DeleteMessage(channel, timestamp string) error {
	client := s.bot.Client()

	_, _, err := client.DeleteMessage(channel, timestamp)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to delete message - %s\n", err)
		return err
	}
	utils.InfoLogger.Printf("Message successfully deleted in channel(%s) from (%s)\n", channel, timestamp)
	return nil
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
