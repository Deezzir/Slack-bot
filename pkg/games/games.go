package games

import (
	"slack-bot/pkg/utils"
	"time"
)

type Game interface {
	init()
	String() (string, error)
	Start(id string)
	IsStarted() bool
	GetTimer() time.Duration
	watcher()
}

var (
	tictactoe *utils.Singleton
)

const Timeout time.Duration = 10 * time.Minute

func resetTicTacToe() {
	tictactoe = nil
}
