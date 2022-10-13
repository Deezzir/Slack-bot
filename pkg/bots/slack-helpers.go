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

func getUser(client *slack.Client, event *slacker.MessageEvent) (*slack.User, error) {
	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func getBotID(botCtx slacker.BotContext) string {
	botID := utils.ExtractTxt(utils.MentionRegex, botCtx.Event().Text)
	if botID == "" {
		botID = "noxu-bot"
	}

	return botID
}

func deleteTicTacToeMsg(botCtx slacker.BotContext, timestamp string) error {
	err := slackBot.Instance.(*SlackBot).DeleteMessage(botCtx.Event().Channel, timestamp)
	if err != nil {
		return err
	}
	return nil
}

func getTicTacToeString(game *games.TicTacToe) string {
	r := fmt.Sprintf("Time left: *%s*\n\n", game.GetTimer())
	r += fmt.Sprintf("<@%s>: %s\n", game.GetBotID(), game.GetAISymbol())
	r += fmt.Sprintf("<@%s>: %s\n", game.GetUserID(), game.GetPlayerSymbol())
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
			ok, winner := g.GetWinner()
			if !ok {
				utils.InfoLogger.Println("Tic-Tac-Toe game was stopped")
				tictactoeTimestamp = ""
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
			if tictactoeTimestamp != "" {
				deleteTicTacToeMsg(botCtx, tictactoeTimestamp)
			}
			break
		}
		if g.GetTimer() <= 0 {
			utils.InfoLogger.Println("Tic-Tac-Toe game has timed out")
			g.Stop()
			msg = fmt.Sprintf("Tick-Tack-Toe game with <@%s> has finished. Timed out.\n", g.GetUserID())
			msg += "Use `tictactoe start` to start a new game\n"

			break
		}
		time.Sleep(1 * time.Second)
		g.SetTimer(g.GetTimer() - 1*time.Second)
	}

	_, err := slackBot.Instance.(*SlackBot).PostMessage(botCtx.Event().Channel, "", msg)
	if err != nil {
		utils.ErrorLogger.Printf("Error while posting message - %s\n", err.Error())
	}
}
