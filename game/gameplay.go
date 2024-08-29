package game

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math"
	"sync"
	"time"
)

type Game struct {
	// Id game id
	Id          string
	Player1Id   string
	Player1Name string
	Player2Id   string
	Player2Name string

	Board      [3][3]int
	State      State
	runningSum RunningSum

	mu *sync.Mutex
}

// RunningSum storing and incrementing the sum of rows, cols, diagonals for each move,
// we know there is a win if one of the sum grows to 3
type RunningSum struct {
	rowSum             [3]int
	columnSum          [3]int
	diagonalSum        int
	reverseDiagonalSum int
}

type State struct {
	Player2Turn bool
	End         bool
	Player1Won  bool
	Player2Won  bool
	Draw        bool
	EndTime     time.Time
}

// predefined errors

var (
	GameIdNotfoundErr  = errors.New("game id not found")
	InvalidPlayerIdErr = errors.New("invalid player id")
)

type NewGameFactory struct{}

func (gf *NewGameFactory) CreateGame(playerName string) *Game {
	return &Game{
		Id:          uuid.NewString(),
		Player1Id:   uuid.NewString(),
		Player1Name: playerName,
		mu:          &sync.Mutex{},
	}
}

// Join player can join a game
func (g *Game) Join(gameId string, playerId string, playerName string) error {
	if gameId != g.Id {
		return GameIdNotfoundErr
	}
	if playerId == "" {
		return InvalidPlayerIdErr
	}
	if playerId == g.Player1Id || playerId == g.Player2Id {
		return errors.New("you have already join this game. please play a move")
	}
	if g.Player2Id != "" {
		return errors.New("game is filled with 2 players already. try joining another game")
	}

	g.Player2Id = playerId
	g.Player2Name = playerName
	return nil
}

// Move makes a move on the game board and record the state of the game
// Player 1's move will be represented by 1 and player 2's move will be represented by -1, empty slot represented by 0
func (g *Game) Move(playerId string, row int, col int) error {
	// size of board
	var n int = 3
	g.mu.Lock()
	defer g.mu.Unlock()
	// state check
	if g.State.End {
		return errors.New("game finished")
	}
	if playerId != g.Player1Id && playerId != g.Player2Id {
		return InvalidPlayerIdErr
	}
	// player turn check
	if (g.State.Player2Turn && playerId != g.Player2Id) || (!g.State.Player2Turn && playerId != g.Player1Id) {
		return errors.New("invalid move. Please wait for other player to move")
	}

	// boundary check
	if row < 0 || col < 0 || row >= n || col >= n {
		return errors.New("invalid move. constraints: 0 <= row < 3, 0 <= column < 3")
	}
	// position filled check
	if g.Board[row][col] != 0 {
		return errors.New("invalid move. position is filled")
	}

	// fill position for current move
	move := 1
	if g.State.Player2Turn {
		move = -1
	}
	g.Board[row][col] = move
	g.runningSum.rowSum[row] += move
	g.runningSum.columnSum[col] += move
	if row == col {
		g.runningSum.diagonalSum += move
	}
	// reverse diagonal position pattern row = n - 1- column
	if row == n-1-col {
		g.runningSum.reverseDiagonalSum += move
	}

	// check for wins
	if math.Abs(float64(g.runningSum.rowSum[row])) == float64(n) ||
		math.Abs(float64(g.runningSum.rowSum[row])) == float64(n) ||
		math.Abs(float64(g.runningSum.rowSum[row])) == float64(n) ||
		math.Abs(float64(g.runningSum.rowSum[row])) == float64(n) {
		g.State.End = true
		if g.State.Player2Turn {
			g.State.Player2Won = true
		} else {
			g.State.Player1Won = true
		}
		g.State.EndTime = time.Now()
		return nil
	}

	// check if it is a draw - all position filled
	var emptySlot bool
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if g.Board[i][j] == 0 {
				emptySlot = true
			}
		}
	}
	if !emptySlot {
		g.State.End = true
		g.State.Draw = true
		g.State.EndTime = time.Now()
		return nil
	}

	// game keep going, switch handle to the next player
	g.State.Player2Turn = !g.State.Player2Turn
	return nil
}

func (g *Game) ShowGameState(sessionId string) string {
	lineHeader := fmt.Sprintf("Session: %s. Player1: %s represent by X. Player2: %s represent by O.", sessionId, g.Player1Name, g.Player2Name)
	lineBoard := ""
	// m takes value 1, -1, or 0
	for _, row := range g.Board {
		lineBoard += "\n___________\n"
		for _, col := range row {
			if col == 1 {
				lineBoard += "X"
			} else if col == -1 {
				lineBoard += "O"
			} else {
				lineBoard += " "
			}
			lineBoard += " | "
		}
	}
	lineBoard += "\n___________\n"

	var lineState = "Game state: "
	if g.State.Draw {
		lineState += fmt.Sprintf("Draw")
	} else if g.State.Player1Won {
		lineState += fmt.Sprintf("%s Won", g.Player1Name)
	} else if g.State.Player2Won {
		lineState += fmt.Sprintf("%s Won", g.Player2Name)
	} else if g.State.Player2Turn {
		lineState += fmt.Sprintf("%s Turn", g.Player2Name)
	} else {
		lineState += fmt.Sprintf("%s Turn", g.Player1Name)
	}

	x := fmt.Sprintf("%s\n%s\n%s", lineHeader, lineBoard, lineState)
	fmt.Println(x)
	return fmt.Sprintf("%s\n%s\n%s", lineHeader, lineBoard, lineState)
}

// EndGame end the game
func (g *Game) EndGame(gameId string, playerId string) error {
	if gameId != g.Id {
		return GameIdNotfoundErr
	}
	if playerId == "" || (playerId != g.Player1Id && playerId != g.Player2Id) {
		return InvalidPlayerIdErr
	}
	g.State.End = true
	g.State.EndTime = time.Now()
	// set a tie if no one wins for simplicity
	if !g.State.Player1Won && !g.State.Player2Won {
		g.State.Draw = true
	}
	return nil
}
