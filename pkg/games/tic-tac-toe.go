package games

import (
	"context"
	"fmt"
	"math/rand"
	"slack-bot/pkg/utils"
	"time"
)

type cell uint8

const (
	empty cell = iota
	cross
	circle
)

type TicTacToe struct {
	Board  [3][3]cell
	Player string
	AI     string
}

func (g *TicTacToe) Init() {
	g.Board = [3][3]cell{
		{empty, empty, empty},
		{empty, empty, empty},
		{empty, empty, empty},
	}

	rand.Seed(time.Now().UnixNano())
	ran := rand.Float64()

	if ran > 0.5 {
		g.Player = "X"
		g.AI = "O"
	} else {
		g.Player = "O"
		g.AI = "X"
	}

}

func (g *TicTacToe) String() (string, error) {
	if g.Player == "" || g.AI == "" {
		return "", fmt.Errorf("board is not initialized")
	}

	var board string
	for i, row := range g.Board {
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

func Watcher(ctx context.Context) {
	newCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	for {
		<-newCtx.Done()
		utils.InfoLogger.Println("Tic-Tac-Toe game has ended")
		TicTacToeCtx = nil
		return
	}
}
