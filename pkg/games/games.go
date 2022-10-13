package games

import (
	"time"
)

type Game interface {
	init()
	String() string

	Start(userID, botID string) bool
	Stop()

	IsRunning() bool

	GetTimer() time.Duration
	SetTimer(time.Duration)

	GetWinner() (bool, string)
}

type singleton struct {
	Instance Game
}

var (
	tictactoe *singleton
)

const Timeout time.Duration = 10 * time.Minute

func ResetTicTacToe() {
	tictactoe = nil
}
