package api

import (
	"github.com/minozihao/tic-tac-toe-server/game"
	"testing"
)

func TestSession_CreateGameInSession(t *testing.T) {
	type fields struct {
		Id         string
		ActiveGame *game.Game
	}
	type args struct {
		playerName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *game.Game
		wantErr bool
	}{
		{
			name: "general",
			fields: fields{
				Id:         "abc",
				ActiveGame: nil,
			},
			args: args{
				playerName: "bob",
			},
			want: &game.Game{
				Player1Name: "bob",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Id:         tt.fields.Id,
				ActiveGame: tt.fields.ActiveGame,
			}
			got, err := s.CreateGameInSession(tt.args.playerName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateGameInSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Player1Name != tt.want.Player1Name || got.Player2Name != tt.want.Player2Name {
				t.Error("expect same name")
			}
		})
	}
}

func TestSession_GetSessionInfo(t *testing.T) {
	type fields struct {
		Id         string
		ActiveGame *game.Game
	}
	tests := []struct {
		name            string
		fields          fields
		expectSessionId string
		expectGameId    string
	}{
		{
			name: "general",
			fields: fields{
				Id: "test_session_id",
				ActiveGame: &game.Game{
					Id: "test_game_id",
				},
			},
			expectSessionId: "test_session_id",
			expectGameId:    "test_game_id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Id:         tt.fields.Id,
				ActiveGame: tt.fields.ActiveGame,
			}
			got, got1 := s.GetSessionInfo()
			if got != tt.expectSessionId {
				t.Errorf("GetSessionInfo() got = %v, want %v", got, tt.expectSessionId)
			}
			if got1 != tt.expectGameId {
				t.Errorf("GetSessionInfo() got1 = %v, want1 %v", got1, tt.expectGameId)
			}
		})
	}
}

func TestSession_JoinGame(t *testing.T) {
	type fields struct {
		Id         string
		ActiveGame *game.Game
	}
	type args struct {
		gameId     string
		playerName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "general success case",
			fields: fields{
				Id: "test-session-id",
				ActiveGame: &game.Game{
					Id: "test_game_id",
				},
			},
			args: args{
				gameId:     "test_game_id",
				playerName: "john",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Id:         tt.fields.Id,
				ActiveGame: tt.fields.ActiveGame,
			}
			got, err := s.JoinGame(tt.args.gameId, tt.args.playerName)
			if (err != nil) != tt.wantErr {
				t.Errorf("JoinGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Error("JoinGame() expect returns a generated uuid player id for player")
			}
		})
	}
}
