package bots

import (
	"fmt"
	"slack-bot/pkg/games"
	"slack-bot/pkg/utils"
	"time"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

var (
	errorMsg           = "Something went wrong, sorry :("
	tictactoeTimestamp = ""
)

func getUser(client *slack.Client, event *slacker.MessageEvent) (*slack.User, bool) {
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get user info: %s\n", err)
		return nil, false
	}
	return user, true
}

func getBotID(botCtx slacker.BotContext) string {
	botID := utils.ExtractTxt(utils.MentionRegex, botCtx.Event().Text)
	if botID == "" {
		botID = "noxu-bot"
	}
	return botID
}

func deleteMsg(botCtx slacker.BotContext, timestamp string) bool {
	err := slackBot.Instance.(*SlackBot).DeleteMessage(botCtx.Event().Channel, timestamp)
	return err == nil
}

func deleteTicTacToeMsg(botCtx slacker.BotContext) {
	if tictactoeTimestamp != "" {
		deleteMsg(botCtx, tictactoeTimestamp)
		tictactoeTimestamp = ""
	}
}

func postMsg(botCtx slacker.BotContext, pretext, text string) (string, bool) {
	timestamp, err := slackBot.Instance.(*SlackBot).PostMessage(botCtx.Event().Channel, pretext, text)
	if err != nil {
		return "", false
	}
	return timestamp, true
}

func getTicTacToeString(game *games.TicTacToe) string {
	r := fmt.Sprintf("Time left: *%s*. Difficulty: `%s`\n\n", game.GetTimer(), game.GetLevel())
	r += fmt.Sprintf("<@%s>: %s\n", game.GetBotID(), game.GetBotSymbol())
	r += fmt.Sprintf("<@%s>: %s\n", game.GetUserID(), game.GetUserSymbol())
	r += fmt.Sprintf("```%s```\n", game.String())

	return r
}

func watchTicTacToe(botCtx slacker.BotContext, g *games.TicTacToe, timeout time.Duration) {
	if g == nil || !g.IsRunning() {
		utils.ErrorLogger.Println("Tic-Tac-Toe game is not started")
		return
	}
	defer games.ResetTicTacToe()

	var msg string
	for {
		if !g.IsRunning() {
			ok, winner := g.GetWinnerID()
			if !ok {
				utils.InfoLogger.Println("Tic-Tac-Toe game was stopped")
				return
			}
			if winner == "" {
				utils.InfoLogger.Println("Tic-Tac-Toe game was a draw")
				msg = "Tic-Tac-Toe game was a draw\n"
			} else {
				utils.InfoLogger.Printf("Tic-Tac-Toe game was won by %s\n", winner)
				msg = fmt.Sprintf("Tic-Tac-Toe game was won by <@%s>\n", winner)
			}

			msg += getTicTacToeString(g)
			deleteTicTacToeMsg(botCtx)
			break
		}
		if g.GetTimer() <= 0 {
			utils.InfoLogger.Println("Tic-Tac-Toe game has timed out")

			msg = fmt.Sprintf("Tick-Tack-Toe game with <@%s> has finished. Timed out.\n", g.GetUserID())
			msg += "Use `tictactoe start` to start a new game\n"
			g.Stop()

			break
		}
		time.Sleep(1 * time.Second)
		g.SetTimer(g.GetTimer() - 1*time.Second)
	}

	_, ok := postMsg(botCtx, "", msg)
	if !ok {
		utils.ErrorLogger.Println("Failed to post Tic-Tac-Toe game message")
	}
}
