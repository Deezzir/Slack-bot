package games

import (
	"context"
)

type Game interface {
	Init()
	String() (string, error)
	Watcher(ctx context.Context)
}

type ctxKey int

const (
	CtxKeyTicTacToe ctxKey = iota
)

var TicTacToeCtx context.Context
