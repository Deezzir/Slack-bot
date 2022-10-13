package games

import (
	"fmt"
	"math"
	"math/rand"
	"slack-bot/pkg/utils"
	"time"
)

type cell uint8
type symbol string
type final uint8
type level uint8

const (
	maxDepth uint8 = 9

	easy       level = 16
	medium     level = 8
	hard       level = 6
	impossible level = 4
)

const (
	empty cell = 1 + iota
	circle
	cross
)

const (
	none final = iota
	userWin
	tie
	botWin
)

const (
	circleSymbol symbol = "O"
	crossSymbol  symbol = "X"
)

type player struct {
	symbol  symbol
	cell    cell
	ID      string
	current bool
}

type TicTacToe struct {
	user   player
	bot    player
	winner final

	board         [3][3]cell
	openCellCount uint8

	level level

	running bool
	timer   time.Duration
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
		g.user.symbol = crossSymbol
		g.user.cell = cross
		g.user.current = true

		g.bot.symbol = circleSymbol
		g.bot.cell = circle
		g.bot.current = false
	} else {
		g.user.symbol = circleSymbol
		g.user.cell = circle
		g.user.current = false

		g.bot.symbol = crossSymbol
		g.bot.cell = cross
		g.bot.current = true
	}

	g.winner = none
	g.level = medium
	g.running = false
	g.timer = Timeout
}

// **************************************
// Public methods
// **************************************

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
				board += " "
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

func (g *TicTacToe) Play(pos string) error {
	if !g.running {
		return fmt.Errorf("tic-Tac-Toe game is not running")
	}

	err := g.playUser(pos)
	if err != nil {
		return err
	}

	g.playBot()

	return nil
}

func (g *TicTacToe) Start(userID, botID string) bool {
	if g.running {
		utils.ErrorLogger.Printf("Failed to start game: already running")
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

func (g *TicTacToe) GetWinnerID() (bool, string) {
	switch g.winner {
	case none:
		return false, ""
	case tie:
		return true, ""
	case userWin:
		return true, g.user.ID
	case botWin:
		return true, g.bot.ID
	default:
		return false, ""
	}
}

func (g *TicTacToe) GetUserSymbol() string {
	return string(g.user.symbol)
}

func (g *TicTacToe) GetBotSymbol() string {
	return string(g.bot.symbol)
}

func (g *TicTacToe) GetTimer() time.Duration {
	return g.timer
}

func (g *TicTacToe) SetTimer(t time.Duration) {
	g.timer = t
}

func (g *TicTacToe) GetLevel() string {
	switch g.level {
	case easy:
		return "easy"
	case medium:
		return "medium"
	case hard:
		return "hard"
	case impossible:
		return "impossible"
	default:
		return "unknown"
	}
}

func (g *TicTacToe) SetLevel(l string) error {
	switch l {
	case "easy":
		g.level = easy
	case "medium":
		g.level = medium
	case "hard":
		g.level = hard
	case "impossible":
		g.level = impossible
	default:
		return fmt.Errorf("invalid level")
	}
	return nil
}

func (g *TicTacToe) IsRunning() bool {
	return g.running
}

// **************************************
// Private methods
// **************************************

func (g *TicTacToe) getWinner() final {
	var winnner final = none

	//Check rows and cols
	for i := 0; i < 3; i++ {
		if g.board[i][0] != empty && g.board[i][0] == g.board[i][1] && g.board[i][2] == g.board[i][0] {
			winnner = g.getWinnnerByCell(g.board[i][0])
		}
		if g.board[0][i] != empty && g.board[0][i] == g.board[1][i] && g.board[2][i] == g.board[0][i] {
			winnner = g.getWinnnerByCell(g.board[0][i])
		}
	}

	//Check diagonals
	if g.board[0][0] != empty && g.board[0][0] == g.board[1][1] && g.board[2][2] == g.board[0][0] {
		winnner = g.getWinnnerByCell(g.board[0][0])
	}
	if g.board[0][2] != empty && g.board[0][2] == g.board[1][1] && g.board[0][2] == g.board[2][0] {
		winnner = g.getWinnnerByCell(g.board[0][2])
	}

	//Check for tie
	if g.openCellCount == 0 && winnner == none {
		winnner = tie
	}

	return winnner
}

func (g *TicTacToe) getWinnnerByCell(cell cell) final {
	if cell == g.user.cell {
		return userWin
	} else if cell == g.bot.cell {
		return botWin
	}
	return none
}

func (g *TicTacToe) setWinner(winner final) {
	if winner != none && g.running {
		g.Stop()
		g.winner = winner
	}
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

		g.board[row][col] = g.user.cell
		g.user.current = false
		g.bot.current = true
		g.openCellCount--

		winner := g.getWinner()
		g.setWinner(winner)
	}
	return nil
}

func (g *TicTacToe) playBot() error {
	if g.bot.current && g.running && g.openCellCount != 0 {
		var bestScore int8 = math.MinInt8
		bestRow := -1
		bestCol := -1

		for row := 0; row < 3; row++ {
			for col := 0; col < 3; col++ {
				if g.board[row][col] != empty {
					continue
				}
				g.board[row][col] = g.bot.cell
				g.openCellCount--

				score := minimax(g, 0, false)

				g.board[row][col] = empty
				g.openCellCount++

				if score > bestScore {
					bestScore = score
					bestRow = row
					bestCol = col
				}
			}
		}

		g.board[bestRow][bestCol] = g.bot.cell
		g.bot.current = false
		g.user.current = true
		g.openCellCount--

		winner := g.getWinner()
		g.setWinner(winner)
	}
	return nil
}

func minimax(g *TicTacToe, depth uint8, max bool) int8 {
	winner := g.getWinner()
	if winner != none {
		return int8(winner)
	}

	var bestScore int8
	var playerCell cell
	var eval func(x, y int8) int8

	if max {
		bestScore = math.MinInt8
		eval = utils.Max
		playerCell = g.bot.cell
	} else {
		bestScore = math.MaxInt8
		eval = utils.Min
		playerCell = g.user.cell
	}

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if g.board[row][col] != empty {
				continue
			}
			g.board[row][col] = playerCell
			g.openCellCount--

			bestScore = eval(bestScore, minimax(g, depth+1, !max))

			g.board[row][col] = empty
			g.openCellCount++
		}
	}
	return bestScore
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
