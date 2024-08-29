package game

import (
	"errors"
	"reflect"
	"sync"
	"testing"
)

func TestGame_EndGame(t *testing.T) {
	type fields struct {
		Id          string
		Player1Id   string
		Player1Name string
		Player2Id   string
		Player2Name string
		Board       [3][3]int
		State       State
		runningSum  RunningSum
		mu          *sync.Mutex
	}
	type args struct {
		gameId   string
		playerId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "general success",
			fields: fields{
				Id:        "test_game_id",
				Player1Id: "test_player1_id",
			},
			args: args{
				gameId:   "test_game_id",
				playerId: "test_player1_id",
			},
			wantErr: false,
		},
		{
			name: "err invalid player id",
			fields: fields{
				Id:        "test_game_id",
				Player1Id: "test_player1_id",
			},
			args: args{
				gameId:   "test_game_id",
				playerId: "invalid player id",
			},
			wantErr: true,
		},
		{
			name: "err invalid game id",
			fields: fields{
				Id:        "test_game_id",
				Player1Id: "test_player1_id",
			},
			args: args{
				gameId:   "invalid_game_id",
				playerId: "test_player1_id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Id:          tt.fields.Id,
				Player1Id:   tt.fields.Player1Id,
				Player1Name: tt.fields.Player1Name,
				Player2Id:   tt.fields.Player2Id,
				Player2Name: tt.fields.Player2Name,
				Board:       tt.fields.Board,
				State:       tt.fields.State,
				runningSum:  tt.fields.runningSum,
				mu:          tt.fields.mu,
			}
			if err := g.EndGame(tt.args.gameId, tt.args.playerId); (err != nil) != tt.wantErr {
				t.Errorf("EndGame() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGame_Join(t *testing.T) {
	type fields struct {
		Id          string
		Player1Id   string
		Player1Name string
		Player2Id   string
		Player2Name string
		Board       [3][3]int
		State       State
		runningSum  RunningSum
		mu          *sync.Mutex
	}
	type args struct {
		gameId     string
		playerId   string
		playerName string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantErrType error
	}{
		{
			name: "general success",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
			},
			args: args{
				gameId:     "test_game_id",
				playerId:   "test_player2_id",
				playerName: "john",
			},
			wantErr:     false,
			wantErrType: nil,
		},
		{
			name: "err duplicate player name",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
			},
			args: args{
				gameId:     "test_game_id",
				playerId:   "test_player2_id",
				playerName: "bob",
			},
			wantErr:     true,
			wantErrType: DuplicatePlayerNameErr,
		},
		{
			name: "GameFilledWithMaxPlayerErr",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
			},
			args: args{
				gameId:     "test_game_id",
				playerId:   "test_player3_id",
				playerName: "charlie",
			},
			wantErr:     true,
			wantErrType: GameFilledWithMaxPlayerErr,
		},
		{
			name: "AlreadyJoinGameErr",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
			},
			args: args{
				gameId:     "test_game_id",
				playerId:   "test_player2_id",
				playerName: "john",
			},
			wantErr:     true,
			wantErrType: AlreadyJoinGameErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Id:          tt.fields.Id,
				Player1Id:   tt.fields.Player1Id,
				Player1Name: tt.fields.Player1Name,
				Player2Id:   tt.fields.Player2Id,
				Player2Name: tt.fields.Player2Name,
				Board:       tt.fields.Board,
				State:       tt.fields.State,
				runningSum:  tt.fields.runningSum,
				mu:          tt.fields.mu,
			}
			err := g.Join(tt.args.gameId, tt.args.playerId, tt.args.playerName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Join() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !errors.Is(err, tt.wantErrType) {
				t.Errorf("Join() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGame_Move_FailCases(t *testing.T) {
	type fields struct {
		Id          string
		Player1Id   string
		Player1Name string
		Player2Id   string
		Player2Name string
		Board       [3][3]int
		State       State
		runningSum  RunningSum
		mu          *sync.Mutex
	}
	type args struct {
		playerId string
		row      int
		col      int
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantErrType error
	}{
		{
			name: "GameAlreadyFinishedErr",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				State: State{
					End: true,
				},
				mu: &sync.Mutex{},
			},
			args: args{
				playerId: "test_player1_id",
				row:      0,
				col:      0,
			},
			wantErr:     true,
			wantErrType: GameAlreadyFinishedErr,
		},
		{
			name: "InvalidPlayerIdErr",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				mu:          &sync.Mutex{},
			},
			args: args{
				playerId: "test_player3_id",
				row:      0,
				col:      0,
			},
			wantErr:     true,
			wantErrType: InvalidPlayerIdErr,
		},
		{
			name: "AnotherPlayerMoveTurnErr",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				State: State{
					Player2Turn: false,
				},
				mu: &sync.Mutex{},
			},
			args: args{
				playerId: "test_player2_id",
				row:      0,
				col:      0,
			},
			wantErr:     true,
			wantErrType: AnotherPlayerMoveTurnErr,
		},
		{
			name: "AnotherPlayerMoveTurnErr2",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				State: State{
					Player2Turn: true,
				},
				mu: &sync.Mutex{},
			},
			args: args{
				playerId: "test_player1_id",
				row:      0,
				col:      0,
			},
			wantErr:     true,
			wantErrType: AnotherPlayerMoveTurnErr,
		},
		{
			name: "InvalidMoveErrNegativeRow",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				mu:          &sync.Mutex{},
			},
			args: args{
				playerId: "test_player1_id",
				row:      -1,
				col:      0,
			},
			wantErr:     true,
			wantErrType: InvalidMoveErr,
		},
		{
			name: "InvalidMoveErrNegativeCol",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				mu:          &sync.Mutex{},
			},
			args: args{
				playerId: "test_player1_id",
				row:      0,
				col:      -1,
			},
			wantErr:     true,
			wantErrType: InvalidMoveErr,
		},
		{
			name: "InvalidMoveErrRowExceedsBoundary",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				mu:          &sync.Mutex{},
			},
			args: args{
				playerId: "test_player1_id",
				row:      3,
				col:      0,
			},
			wantErr:     true,
			wantErrType: InvalidMoveErr,
		},
		{
			name: "InvalidMoveErrColumnExceedsBoundary",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				mu:          &sync.Mutex{},
			},
			args: args{
				playerId: "test_player1_id",
				row:      0,
				col:      3,
			},
			wantErr:     true,
			wantErrType: InvalidMoveErr,
		}, {
			name: "MovePositionFilledErr",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{1, -1, 1}, {1, -1, 0}, {1, -1, 1},
				},
				mu: &sync.Mutex{},
			},
			args: args{
				playerId: "test_player1_id",
				row:      0,
				col:      0,
			},
			wantErr:     true,
			wantErrType: MovePositionFilledErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Id:          tt.fields.Id,
				Player1Id:   tt.fields.Player1Id,
				Player1Name: tt.fields.Player1Name,
				Player2Id:   tt.fields.Player2Id,
				Player2Name: tt.fields.Player2Name,
				Board:       tt.fields.Board,
				State:       tt.fields.State,
				runningSum:  tt.fields.runningSum,
				mu:          tt.fields.mu,
			}
			err := g.Move(tt.args.playerId, tt.args.row, tt.args.col)
			if (err != nil) != tt.wantErr {
				t.Errorf("Move() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !errors.Is(err, tt.wantErrType) {
				t.Errorf("Join() error = %v, wantErr %v", err, tt.wantErrType)
			}
		})
	}
}

func TestGame_GeneralMove(t *testing.T) {
	type fields struct {
		Id          string
		Player1Id   string
		Player1Name string
		Player2Id   string
		Player2Name string
		Board       [3][3]int
		State       State
		runningSum  RunningSum
		mu          *sync.Mutex
	}
	type args struct {
		playerId string
		row      int
		col      int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    *Game
	}{
		{
			name: "General Move turn switch to another user",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{0, 0, 1}, {1, -1, 0}, {1, -1, 1},
				},
				mu: &sync.Mutex{},
			},
			args: args{
				playerId: "test_player1_id",
				row:      0,
				col:      0,
			},
			wantErr: false,
			want: &Game{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{1, 0, 1}, {1, -1, 0}, {1, -1, 1},
				},
				State: State{
					Player2Turn: true,
					End:         false,
					Player1Won:  false,
					Player2Won:  false,
					Draw:        false,
				},
			},
		},
		{
			name: "Player 1 row win after move",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{0, 1, 1}, {1, -1, 0}, {1, -1, 1},
				},
				runningSum: RunningSum{
					rowSum: [3]int{2, 0, 0},
				},
				mu: &sync.Mutex{},
			},
			args: args{
				playerId: "test_player1_id",
				row:      0,
				col:      0,
			},
			wantErr: false,
			want: &Game{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{1, 1, 1}, {1, -1, 0}, {1, -1, 1},
				},
				State: State{
					Player2Turn: false,
					End:         true,
					Player1Won:  true,
					Player2Won:  false,
					Draw:        false,
				},
			},
		},
		{
			name: "Player 1 colum win after move",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{0, 0, 1}, {1, -1, 0}, {1, -1, 1},
				},
				runningSum: RunningSum{
					rowSum:    [3]int{0, 0, 0},
					columnSum: [3]int{2, 0, 0},
				},
				mu: &sync.Mutex{},
			},
			args: args{
				playerId: "test_player1_id",
				row:      0,
				col:      0,
			},
			wantErr: false,
			want: &Game{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{1, 0, 1}, {1, -1, 0}, {1, -1, 1},
				},
				State: State{
					Player2Turn: false,
					End:         true,
					Player1Won:  true,
					Player2Won:  false,
					Draw:        false,
				},
			},
		},
		{
			name: "Player 1 diagonal win after move",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{0, 0, 1}, {-1, 1, 0}, {0, -1, 1},
				},
				runningSum: RunningSum{
					rowSum:      [3]int{0, 0, 0},
					columnSum:   [3]int{0, 0, 0},
					diagonalSum: 2,
				},
				mu: &sync.Mutex{},
			},
			args: args{
				playerId: "test_player1_id",
				row:      0,
				col:      0,
			},
			wantErr: false,
			want: &Game{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{1, 0, 1}, {-1, 1, 0}, {0, -1, 1},
				},
				State: State{
					Player2Turn: false,
					End:         true,
					Player1Won:  true,
					Player2Won:  false,
					Draw:        false,
				},
			},
		},
		{
			name: "Player 1 reverse diagonal win after move",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{-1, 0, 0}, {-1, 1, 0}, {1, -1, 1},
				},
				runningSum: RunningSum{
					rowSum:             [3]int{0, 0, 0},
					columnSum:          [3]int{0, 0, 0},
					diagonalSum:        0,
					reverseDiagonalSum: 2,
				},
				mu: &sync.Mutex{},
			},
			args: args{
				playerId: "test_player1_id",
				row:      0,
				col:      2,
			},
			wantErr: false,
			want: &Game{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{-1, 0, 1}, {-1, 1, 0}, {1, -1, 1},
				},
				State: State{
					Player2Turn: false,
					End:         true,
					Player1Won:  true,
					Player2Won:  false,
					Draw:        false,
				},
			},
		},
		{
			name: "Draw when all position filled",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{-1, 1, 0}, {-1, 1, 1}, {1, -1, -1},
				},
				mu: &sync.Mutex{},
			},
			args: args{
				playerId: "test_player1_id",
				row:      0,
				col:      2,
			},
			wantErr: false,
			want: &Game{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{-1, 1, 1}, {-1, 1, 1}, {1, -1, -1},
				},
				State: State{
					Player2Turn: false,
					End:         true,
					Player1Won:  false,
					Player2Won:  false,
					Draw:        true,
				},
			},
		},
		{
			name: "Player 2 row win",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{-1, -1, 0}, {-1, 0, 0}, {1, 0, 0},
				},
				runningSum: RunningSum{
					rowSum:             [3]int{-2, 0, 0},
					columnSum:          [3]int{0, 0, 0},
					diagonalSum:        0,
					reverseDiagonalSum: 0,
				},
				State: State{Player2Turn: true},
				mu:    &sync.Mutex{},
			},
			args: args{
				playerId: "test_player2_id",
				row:      0,
				col:      2,
			},
			wantErr: false,
			want: &Game{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{-1, -1, -1}, {-1, 0, 0}, {1, 0, 0},
				},
				State: State{
					Player2Turn: true,
					End:         true,
					Player1Won:  false,
					Player2Won:  true,
					Draw:        false,
				},
			},
		},
		{
			name: "Player 2 column win",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{1, 0, -1}, {0, 0, -1}, {1, 0, 0},
				},
				runningSum: RunningSum{
					rowSum:             [3]int{0, 0, 0},
					columnSum:          [3]int{0, 0, -2},
					diagonalSum:        0,
					reverseDiagonalSum: 0,
				},
				State: State{Player2Turn: true},
				mu:    &sync.Mutex{},
			},
			args: args{
				playerId: "test_player2_id",
				row:      2,
				col:      2,
			},
			wantErr: false,
			want: &Game{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{1, 0, -1}, {0, 0, -1}, {1, 0, -1},
				},
				State: State{
					Player2Turn: true,
					End:         true,
					Player1Won:  false,
					Player2Won:  true,
					Draw:        false,
				},
			},
		},
		{
			name: "Player 2 diagonal win",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{-1, 0, 1}, {0, -1, 1}, {1, 0, 0},
				},
				runningSum: RunningSum{
					rowSum:             [3]int{0, 0, 0},
					columnSum:          [3]int{0, 0, 0},
					diagonalSum:        -2,
					reverseDiagonalSum: 0,
				},
				State: State{Player2Turn: true},
				mu:    &sync.Mutex{},
			},
			args: args{
				playerId: "test_player2_id",
				row:      2,
				col:      2,
			},
			wantErr: false,
			want: &Game{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{-1, 0, 1}, {0, -1, 1}, {1, 0, -1},
				},
				State: State{
					Player2Turn: true,
					End:         true,
					Player1Won:  false,
					Player2Won:  true,
					Draw:        false,
				},
			},
		},
		{
			name: "Player 2 reverse diagonal win",
			fields: fields{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{1, 0, -1}, {0, -1, 1}, {0, 0, 1},
				},
				runningSum: RunningSum{
					rowSum:             [3]int{0, 0, 0},
					columnSum:          [3]int{0, 0, 0},
					diagonalSum:        0,
					reverseDiagonalSum: -2,
				},
				State: State{Player2Turn: true},
				mu:    &sync.Mutex{},
			},
			args: args{
				playerId: "test_player2_id",
				row:      2,
				col:      0,
			},
			wantErr: false,
			want: &Game{
				Id:          "test_game_id",
				Player1Id:   "test_player1_id",
				Player1Name: "bob",
				Player2Id:   "test_player2_id",
				Player2Name: "john",
				Board: [3][3]int{
					{1, 0, -1}, {0, -1, 1}, {-1, 0, 1},
				},
				State: State{
					Player2Turn: true,
					End:         true,
					Player1Won:  false,
					Player2Won:  true,
					Draw:        false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				Id:          tt.fields.Id,
				Player1Id:   tt.fields.Player1Id,
				Player1Name: tt.fields.Player1Name,
				Player2Id:   tt.fields.Player2Id,
				Player2Name: tt.fields.Player2Name,
				Board:       tt.fields.Board,
				State:       tt.fields.State,
				runningSum:  tt.fields.runningSum,
				mu:          tt.fields.mu,
			}
			err := g.Move(tt.args.playerId, tt.args.row, tt.args.col)
			if (err != nil) != tt.wantErr {
				t.Errorf("Move() error = %v, wantErr %v", err, tt.wantErr)
			}
			// set to same time for deep equal
			tt.want.State.EndTime = g.State.EndTime

			if !reflect.DeepEqual(g.State, tt.want.State) {
				t.Errorf("Move() state mismatch actualState = %v, wantState %v", g.State, tt.want.State)
			}
			if !reflect.DeepEqual(g.Board, tt.want.Board) {
				t.Errorf("Move() board mismatch actualBoard = %v, wantBoard %v", g.Board, tt.want.Board)
			}
		})
	}
}

func TestNewGameFactory_CreateGame(t *testing.T) {
	type args struct {
		playerName string
	}
	tests := []struct {
		name string
		args args
		want *Game
	}{
		{
			name: "general",
			args: args{
				playerName: "bob",
			},
			want: &Game{
				Player1Name: "bob",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gf := &NewGameFactory{}
			got := gf.CreateGame(tt.args.playerName)
			if got.Player1Name != tt.want.Player1Name {
				t.Errorf("CreateGame() = %v, want %v", got, tt.want)
			} else if got.Player1Id == "" {
				t.Error("expect non empty player id")
			} else if got.Player2Name != "" || got.Player2Id != "" {
				t.Error("expect empty player2 id and name")
			} else if !reflect.DeepEqual(got.State, tt.want.State) {
				t.Errorf("got state %v, want %v", got.State, tt.want.State)
			}
		})
	}
}
