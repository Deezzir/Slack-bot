package games

import (
	"fmt"
	"math/rand"
	"slack-bot/pkg/utils"
	"time"
)

type cell uint8
type playerType uint8

const (
	empty cell = 1 + iota
	circle
	cross
)

const (
	tie playerType = 1 + iota
	circlePlayer
	crossPlayer
)

type player struct {
	playerType playerType
	ID         string
	current    bool
}

type TicTacToe struct {
	user   player
	bot    player
	winner playerType

	board         [3][3]cell
	openCellCount uint8

	running bool
	timer   time.Duration
}

func GetTicTacToe() *singleton {
	if tictactoe == nil {
		utils.GameLock.Lock()
		defer utils.GameLock.Unlock()

		if tictactoe == nil {
			board := &TicTacToe{}
			board.init()

			tictactoe = &singleton{Instance: board}
		}
	}
	return tictactoe
}

func (g *TicTacToe) init() {
	g.board = [3][3]cell{
		{empty, empty, empty},
		{empty, empty, empty},
		{empty, empty, empty},
	}
	g.openCellCount = 9

	rand.Seed(time.Now().UnixNano())
	ran := rand.Float64()

	if ran > 0.5 {
		g.user.playerType = crossPlayer
		g.bot.playerType = circlePlayer
		g.user.current = true
		g.bot.current = false
	} else {
		g.user.playerType = circlePlayer
		g.bot.playerType = crossPlayer
		g.bot.current = true
		g.user.current = false
	}

	g.winner = 0
	g.running = false
	g.timer = Timeout
}

func (g *TicTacToe) checkWinner() {
	ended := false

	//Check rows and cols
	for i := 0; i < 3; i++ {
		if g.board[i][0] != empty && g.board[i][0] == g.board[i][1] && g.board[i][2] == g.board[i][0] {
			ended = true
			g.winner = playerType(g.board[i][0])
		}
		if g.board[0][i] != empty && g.board[0][i] == g.board[1][i] && g.board[2][i] == g.board[0][i] {
			ended = true
			g.winner = playerType(g.board[0][i])
		}
	}

	//Check diagonals
	if g.board[0][0] != empty && g.board[0][0] == g.board[1][1] && g.board[2][2] == g.board[0][0] {
		ended = true
		g.winner = playerType(g.board[0][0])
	}
	if g.board[0][2] != empty && g.board[0][2] == g.board[1][1] && g.board[0][2] == g.board[2][0] {
		ended = true
		g.winner = playerType(g.board[0][2])
	}

	//Check for tie
	if g.openCellCount == 0 && !ended {
		ended = true
		g.winner = tie
	}

	if ended {
		g.Stop()
	}
}

func (g *TicTacToe) Play(pos string) error {
	if !g.running {
		return fmt.Errorf("tic-Tac-Toe game is not running")
	}

	err := g.playUser(pos)
	if err != nil {
		return err
	}
	g.checkWinner()

	err = g.playBot()
	if err != nil {
		return err
	}
	g.checkWinner()

	return nil
}

func parsePosition(pos string) (row, col uint8, ok bool) {
	if len(pos) != 2 {
		return 0, 0, false
	}
	switch pos[0] {
	case 'A', 'a':
		row = 0
	case 'B', 'b':
		row = 1
	case 'C', 'c':
		row = 2
	default:
		return 0, 0, false
	}
	switch pos[1] {
	case '1':
		col = 0
	case '2':
		col = 1
	case '3':
		col = 2
	default:
		return 0, 0, false
	}
	return row, col, true
}

func (g *TicTacToe) playUser(pos string) error {
	if g.running && g.openCellCount != 0 {
		if !g.user.current {
			return fmt.Errorf("not your turn")
		}

		row, col, ok := parsePosition(pos)
		if !ok {
			return fmt.Errorf("position `%s` is invalid", pos)
		}

		if g.board[row][col] != empty {
			return fmt.Errorf("position `%s` is already taken", pos)
		}

		g.board[row][col] = cell(g.user.playerType)
		g.user.current = false
		g.bot.current = true
		g.openCellCount--
	}
	return nil
}

func (g *TicTacToe) playBot() error {
	if g.bot.current && g.running && g.openCellCount != 0 {
		for {
			rand.Seed(time.Now().UnixNano())
			row := rand.Intn(3)
			col := rand.Intn(3)
			if g.board[row][col] != empty {
				continue
			}
			g.board[row][col] = cell(g.bot.playerType)
			g.bot.current = false
			g.user.current = true
			g.openCellCount--
			break
		}
	}
	return nil
}

func (g *TicTacToe) String() string {
	if g.board[0][0] == 0 {
		return "Ooops, its seems like a non initialized game"
	}

	var board string
	for i, row := range g.board {
		switch i {
		case 0:
			board += "  ┌───┬───┬───┐\nA"
		case 1:
			board += "  ├───┼───┼───┤\nB"
		case 2:
			board += "  ├───┼───┼───┤\nC"
		}
		board += " | "
		for _, cell := range row {
			switch cell {
			case cross:
				board += "X"
			case circle:
				board += "O"
			case empty:
				board += "-"
			default:
				return "Heh, found a bug"
			}
			board += " | "
		}
		board += "\n"
	}
	board += "  └───┴───┴───┘\n    1   2   3\n"
	return board
}

func (g *TicTacToe) GetPlayerSymbol() string {
	if g.user.playerType == crossPlayer {
		return "X"
	}
	return "O"
}

func (g *TicTacToe) GetAISymbol() string {
	if g.bot.playerType == crossPlayer {
		return "X"
	}
	return "O"
}

func (g *TicTacToe) GetTimer() time.Duration {
	return g.timer
}

func (g *TicTacToe) SetTimer(t time.Duration) {
	g.timer = t
}

func (g *TicTacToe) IsRunning() bool {
	return g.running
}

func (g *TicTacToe) Start(userID, botID string) bool {
	if g.running {
		utils.ErrorLogger.Printf("Failed to start game: Game is already running")
		return false
	}

	g.user.ID = userID
	g.bot.ID = botID

	g.running = true
	g.playBot()

	utils.InfoLogger.Printf("Tic-Tac-Toe game has strated, will timeout at %s\n", time.Now().Add(g.timer).Format(time.UnixDate))
	return true
}

func (g *TicTacToe) Stop() {
	g.running = false
}

func (g *TicTacToe) IsPlayerTurn() bool {
	return g.user.current
}

func (g *TicTacToe) GetUserID() string {
	return g.user.ID
}

func (g *TicTacToe) GetBotID() string {
	return g.bot.ID
}

func (g *TicTacToe) GetWinner() (bool, string) {
	if g.winner == 0 {
		return false, ""
	}

	if g.winner == tie {
		return true, ""
	}

	if g.winner == g.user.playerType {
		return true, g.user.ID
	}

	return true, g.bot.ID
}
