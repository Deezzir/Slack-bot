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

	GetWinnerID() (bool, string)
}

type singleton struct {
	Instance Game
}

var (
	tictactoe *singleton
)

const Timeout time.Duration = 5 * time.Minute

func ResetTicTacToe() {
	tictactoe = nil
}
