package bots

import (
	"encoding/json"
	"slack-bot/pkg/blob"
	"slack-bot/pkg/config"
	"slack-bot/pkg/games"
	"slack-bot/pkg/utils"
	"strings"

	"github.com/krognol/go-wolfram"
	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"

	"fmt"
	"strconv"
	"time"
)

var slackGetAge = SlackCommand{
	Name: "birth year {year}",
	CommandDefinition: &slacker.CommandDefinition{
		Description: "Calculates your age.",
		Examples:    []string{"birth year 1990"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			user, err := getUser(botCtx.Client(), botCtx.Event())
			if err != nil {
				utils.ErrorLogger.Printf("Failed to get user info: %s\n", err)
				response.Reply(errorMsg)
				return
			}

			param := request.Param("year")
			yearStr := utils.ExtractTxt(utils.HyperlinkRegex, param)
			if yearStr == "" {
				yearStr = param
			}

			year, err := strconv.Atoi(yearStr)
			if err != nil || year < 0 {
				r := fmt.Sprintf("'%s' is an invalid year\n", yearStr)
				response.Reply(r)
			} else {
				age := time.Now().Year() - year
				var r string

				if age < 0 {
					r = fmt.Sprintf("<@%s>You are from the future, go away\n", user.ID)
				} else if age == 0 {
					r = fmt.Sprintf("<@%s> Your Age is *%d*, You are born this year, really?\n", user.ID, age)
				} else if age < 18 {
					r = fmt.Sprintf("<@%s> Your Age is *%d*, You are too young\n", user.ID, age)
				} else if age < 22 {
					r = fmt.Sprintf("<@%s> Your Age is *%d*, So fresh\n", user.ID, age)
				} else if age < 100 {
					r = fmt.Sprintf("<@%s> Your Age is *%d*, Too old\n", user.ID, age)
				} else {
					r = fmt.Sprintf("<@%s> Your Age is *%d*, Probably dead\n", user.ID, age)
				}
				response.Reply(r)
			}
		},
	},
}

var slackYouSuck = SlackCommand{
	Name: "you suck",
	CommandDefinition: &slacker.CommandDefinition{
		Description: "You can tell the bot that it sucks. But it will talk back.",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			user, err := getUser(botCtx.Client(), botCtx.Event())
			if err != nil {
				utils.ErrorLogger.Printf("Failed to get user info: %s\n", err)
				response.Reply(errorMsg)
				return
			}

			r := fmt.Sprintf("<@%s> No, you suck!\nI kwon your IP address btw...", user.ID)
			response.Reply(r)
		},
	},
}

var slackListFiles = SlackCommand{
	Name: "list files",
	CommandDefinition: &slacker.CommandDefinition{
		Description: "List files available for download.",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			files := blob.GetBlobFiles(botCtx.Context(), config.CONTAINER)

			if len(files) == 0 {
				response.Reply("No files available for download")
				return
			}

			r := fmt.Sprintln("List of files available for download:")
			for _, file := range files {
				r += fmt.Sprintf("• `%s`", file.Filename)
				if file.Desc != "" {
					r += fmt.Sprintf(" - %s", file.Desc)
				}
				r += "\n"
			}
			response.Reply(r)
		},
	},
}

var slackGetFile = SlackCommand{
	Name: "get file <file>",
	CommandDefinition: &slacker.CommandDefinition{
		Description: "Get available file.",
		Examples:    []string{"get file dog.jpg", "get file doc.pdf"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			param := request.Param("file")
			filename := utils.ExtractTxt(utils.HyperlinkRegex, param)
			if filename == "" {
				filename = param
			}

			client := botCtx.Client()
			event := botCtx.Event()

			if file, ok := blob.GetBlobFile(botCtx.Context(), config.CONTAINER, filename); ok {
				if file.Data != "" {
					if event.Channel != "" {
						params := slack.FileUploadParameters{
							Content:  file.Data,
							Channels: []string{event.Channel},
						}

						r := fmt.Sprintf("Downloading `%s` ...\n", file.Filename)
						client.PostMessage(event.Channel, slack.MsgOptionText(r, false))
						_, err := client.UploadFile(params)
						if err != nil {
							utils.ErrorLogger.Printf("Failed to upload '%s' file to Slack channel\n", filename)
							response.Reply("Sorry, failed to download the file :'(")
						}
					}
				} else {
					response.Reply("File not found, use `list files` for avaliable files")
				}
			} else {
				response.Reply(errorMsg)
			}
		},
	},
}

var slackPutFile = SlackCommand{
	Name: `put file {filename} "<description>"`,
	CommandDefinition: &slacker.CommandDefinition{
		Description: "Save the provided file to Noxu-bot's memory.",
		Examples: []string{
			"put file dog.jpeg",
			`put file "an important document"`,
			`put file flower.png "just a flower, lol"`,
		},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			// TODO
		},
	},
}

var slackValidateEmail = SlackCommand{
	Name: "validate email <email>",
	CommandDefinition: &slacker.CommandDefinition{
		Description: "Check email validity and verify email domain. Does not check if email exists.",
		Examples:    []string{"validate email deezzir@gmail.com"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			param := request.Param("email")
			email := utils.ExtractTxt(utils.HyperlinkRegex, param)
			if email == "" {
				email = param
			}

			local, domain, ok := utils.NormalizeEmail(email)
			if !ok {
				response.Reply("Please provide a valid email")
				return
			}
			r := fmt.Sprintf("*Email*: `%s@%s`\n", local, domain)
			r += fmt.Sprintf("*Domain*: `%s`\n", domain)

			vdDomain := utils.CheckEmailDomain(domain)

			if vdDomain.Valid {
				if len(vdDomain.Addr) > 0 {
					r += fmt.Sprintf("• *Addresses*: `%s`\n", strings.Join(vdDomain.Addr[:], "`, `"))
				}
				r += fmt.Sprintf("• *has MX*: `%t`\n", vdDomain.HasMX)

				if vdDomain.HasSPF {
					r += fmt.Sprintf("• *SPF Record*: `%s`\n", vdDomain.SPFRecord)
				}

				if vdDomain.HasDMARC {
					r += fmt.Sprintf("• *DMARC Record*: `%s`\n", vdDomain.DMARCRecord)
				}

			} else {
				r += fmt.Sprintf("• *Valid*: `%t`\n", vdDomain.Valid)
			}

			response.Reply(r)
		},
	},
}

var slackAskQuestion = SlackCommand{
	Name: "ask <question>",
	CommandDefinition: &slacker.CommandDefinition{
		Description: "Ask Noxu-bot a question.",
		Examples:    []string{"ask what is the meaning of life?", "ask what is the weather like today?"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			param := request.Param("question")
			query := utils.ExtractTxt(utils.HyperlinkRegex, param)
			if query == "" {
				query = param
			}

			msg, _ := config.WitAIClient.Parse(&witai.MessageRequest{
				Query: query,
			})

			obj, err := json.MarshalIndent(msg, "", "    ")
			if err != nil {
				utils.ErrorLogger.Printf("Failed to encode Wit AI response - %s\n", err)
				response.Reply(errorMsg)
				return
			}

			question := gjson.Get(string(obj[:]), "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			res, err := config.WolframClient.GetSpokentAnswerQuery(question.String(), wolfram.Metric, 1000)
			if err != nil {
				utils.ErrorLogger.Printf("Failed to get Wolfram Alpha response - %s\n", err)
				response.Reply(errorMsg)
				return
			}
			if strings.Contains(res, "No spoken result available") || strings.Contains(res, "Wolfram Alpha did not understand your input") {
				response.Reply("Sorry, I don't know the answer to that question")
			} else {
				response.Reply(res)
			}
		},
	},
}

var slackStartTicTacToe = SlackCommand{
	Name: "tictactoe start",
	CommandDefinition: &slacker.CommandDefinition{
		Description: "Start Tic-Tac-Toe game with Noxu-bot",
		Examples:    []string{},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			game := games.GetTicTacToe().Instance.(*games.TicTacToe)

			if !game.IsRunning() {
				user, err := getUser(botCtx.Client(), botCtx.Event())
				if err != nil {
					utils.ErrorLogger.Printf("Failed to get user info: %s\n", err)
					response.Reply(errorMsg)
					return
				}
				botID := getBotID(botCtx)

				if err = game.Start(user.ID, botID); err != nil {
					utils.ErrorLogger.Printf("Failed to start game: %s\n", err)
				}

				if tictactoeTimestamp != "" {
					deleteTicTacToeMsg(botCtx, tictactoeTimestamp)
				}

				r := "Tic-Tac-Toe game started, use `tictactoe play <position>` to play\n"
				r += getTicTacToeString(game)

				go watchTicTacToe(botCtx, game, games.Timeout)

				tictactoeTimestamp, err = slackBot.Instance.(*SlackBot).PostMessage(botCtx.Event().Channel, r, "")
				if err != nil {
					utils.ErrorLogger.Printf("Error while posting message - %s\n", err.Error())
					response.Reply(errorMsg)
				}
			} else {
				r := fmt.Sprintf("Tic-Tac-Toe game already started by <@%s>\n", game.GetUserID())
				r += fmt.Sprintf("Time left: *%s*\n", game.GetTimer())
				response.Reply(r)
			}
		},
	},
}

var slackStopTicTacToe = SlackCommand{
	Name: "tictactoe stop",
	CommandDefinition: &slacker.CommandDefinition{
		Description: "Stop Tic-Tac-Toe game with Noxu-bot",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			game := games.GetTicTacToe().Instance.(*games.TicTacToe)

			if game.IsRunning() {
				game.Stop()
				response.Reply("Tic-Tac-Toe game stopped")
			} else {
				r := "No Tic-Tac-Toe game started.\n"
				response.Reply(r)
			}
		},
	},
}

var slackShowTicTacToe = SlackCommand{
	Name: "tictactoe show",
	CommandDefinition: &slacker.CommandDefinition{
		Description: "Show current Tic-Tac-Toe game with Noxu-bot",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			game := games.GetTicTacToe().Instance.(*games.TicTacToe)

			if game.IsRunning() {
				if tictactoeTimestamp != "" {
					deleteTicTacToeMsg(botCtx, tictactoeTimestamp)
				}

				r := fmt.Sprintf("Tic-Tac-Toe game with <@%s>\n", game.GetUserID())
				r += getTicTacToeString(game)

				var err error
				tictactoeTimestamp, err = slackBot.Instance.(*SlackBot).PostMessage(botCtx.Event().Channel, r, "")
				if err != nil {
					utils.ErrorLogger.Printf("Error while posting message - %s\n", err.Error())
					response.Reply(errorMsg)
				}
			} else {
				r := "No Tic-Tac-Toe game started.\n"
				r += "Use `tictactoe start` to start a new game\n"

				response.Reply(r)
			}
		},
	},
}

var slackPlayTicTacToe = SlackCommand{
	Name: "tictactoe play <position>",
	CommandDefinition: &slacker.CommandDefinition{
		Description: "Play Tic-Tac-Toe game with Noxu-bot",
		Examples:    []string{"tictactoe play A3", "tictactoe play B2", "tictactoe play C1"},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			game := games.GetTicTacToe().Instance.(*games.TicTacToe)

			if game.IsRunning() {
				param := request.Param("position")
				pos := utils.ExtractTxt(utils.HyperlinkRegex, param)
				if pos == "" {
					pos = param
				}

				if game.IsPlayerTurn() {
					err := game.Play(pos)
					if err != nil {
						r := fmt.Sprintf("Tic-Tac-Toe play by <@%s>: %s\n", game.GetUserID(), err)
						response.Reply(r)
						return
					}

					if tictactoeTimestamp != "" {
						deleteTicTacToeMsg(botCtx, tictactoeTimestamp)
					}

					if game.IsRunning() {
						r := fmt.Sprintf("You played Tic-Tac-Toe game with `%s`\n", pos)
						r += getTicTacToeString(game)

						tictactoeTimestamp, err = slackBot.Instance.(*SlackBot).PostMessage(botCtx.Event().Channel, r, "")
						if err != nil {
							utils.ErrorLogger.Printf("Error while posting message - %s\n", err.Error())
							response.Reply(errorMsg)
						}
					}
				} else {
					response.Reply("It's not your turn")
					return
				}
			} else {
				r := "No Tic-Tac-Toe game started.\n"
				r += "Use `tictactoe start` to start a new game\n"

				response.Reply(r)
			}
		},
	},
}
