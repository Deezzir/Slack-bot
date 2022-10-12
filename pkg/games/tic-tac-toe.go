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
	empty cell = iota
	circle
	cross
)

const (
	circlePlayer playerType = 1 + iota
	crossPlayer
)

type TicTacToe struct {
	player  playerType
	bot     playerType
	current playerType

	board   [3][3]cell
	started bool
	timer   time.Duration

	userID string
}

func GetTicTacToe() *utils.Singleton {
	if tictactoe == nil {
		utils.GameLock.Lock()
		defer utils.GameLock.Unlock()

		if tictactoe == nil {
			board := &TicTacToe{}
			board.init()

			tictactoe = &utils.Singleton{Instance: board}
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

	rand.Seed(time.Now().UnixNano())
	ran := rand.Float64()

	if ran > 0.5 {
		g.player = crossPlayer
		g.bot = circlePlayer
		g.current = g.player
	} else {
		g.player = circlePlayer
		g.bot = crossPlayer
		g.current = g.bot
	}

	g.started = false
	g.timer = Timeout
}

func (g *TicTacToe) getWinner() (playerType, bool) {
	//Check rows and cols
	for i := 0; i < 3; i++ {
		if g.board[i][0] != empty && g.board[i][0] == g.board[i][1] && g.board[i][2] == g.board[i][0] {
			return playerType(g.board[i][0]), true
		}
		if g.board[0][i] != empty && g.board[0][i] == g.board[1][i] && g.board[2][i] == g.board[0][i] {
			return playerType(g.board[0][i]), true
		}
	}

	//Check diagonals
	if g.board[0][0] != empty && g.board[0][0] == g.board[1][1] && g.board[2][2] == g.board[0][0] {
		return playerType(g.board[0][0]), true
	}
	if g.board[0][0] != empty && g.board[0][2] == g.board[1][1] && g.board[2][0] == g.board[0][2] {
		return playerType(g.board[0][0]), true
	}

	//Check for tie
	for _, row := range g.board {
		for _, cell := range row {
			if cell == empty {
				return 0, false
			}
		}
	}

	return 0, true // tie
}

func (g *TicTacToe) play(pos string) {

}

func (g *TicTacToe) String() (string, error) {
	if g.player == 0 || g.bot == 0 {
		return "", fmt.Errorf("board is not initialized")
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
				return "", fmt.Errorf("invalid cell value")
			}
			board += " | "
		}
		board += "\n"
	}
	board += "  └───┴───┴───┘\n    1   2   3\n"
	return board, nil
}

func (g *TicTacToe) GetPlayerSymbol() string {
	if g.player == crossPlayer {
		return "X"
	}
	return "O"
}

func (g *TicTacToe) GetAISymbol() string {
	if g.bot == crossPlayer {
		return "X"
	}
	return "O"
}

func (g *TicTacToe) GetTimer() time.Duration {
	return g.timer
}

func (g *TicTacToe) IsStarted() bool {
	return g.started
}

func (g *TicTacToe) Start(id string) {
	g.started = true
	g.userID = id
	go g.watcher()
}

func (g *TicTacToe) GetUserID() string {
	return g.userID
}

func (g *TicTacToe) watcher() {
	defer resetTicTacToe()

	for {
		if g.timer == 0*time.Second {
			utils.InfoLogger.Println("Tic-Tac-Toe game has ended")
			g.started = false
			return
		}
		time.Sleep(1 * time.Second)
		g.timer -= 1 * time.Second
	}
}
