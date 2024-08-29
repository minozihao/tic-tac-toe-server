package game

import (
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
		// TODO: Add test cases.
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
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := g.Join(tt.args.gameId, tt.args.playerId, tt.args.playerName); (err != nil) != tt.wantErr {
				t.Errorf("Join() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGame_Move(t *testing.T) {
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
	}{
		// TODO: Add test cases.
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
			if err := g.Move(tt.args.playerId, tt.args.row, tt.args.col); (err != nil) != tt.wantErr {
				t.Errorf("Move() error = %v, wantErr %v", err, tt.wantErr)
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
