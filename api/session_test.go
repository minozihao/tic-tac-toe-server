package api

import (
	"github.com/gorilla/mux"
	"github.com/minozihao/tic-tac-toe-server/game"
	"github.com/patrickmn/go-cache"
	"reflect"
	"sync"
	"testing"
)

func TestServer_CreateGame(t *testing.T) {
	type fields struct {
		Router        *mux.Router
		Sessions      *sync.Map
		FinishedGames *cache.Cache
	}
	type args struct {
		sessionId  string
		playerName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Router:        tt.fields.Router,
				Sessions:      tt.fields.Sessions,
				FinishedGames: tt.fields.FinishedGames,
			}
			got, got1, err := s.CreateGame(tt.args.sessionId, tt.args.playerName)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateGame() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CreateGame() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestServer_DeleteSession(t *testing.T) {
	type fields struct {
		Router        *mux.Router
		Sessions      *sync.Map
		FinishedGames *cache.Cache
	}
	type args struct {
		sessionId string
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
			s := &Server{
				Router:        tt.fields.Router,
				Sessions:      tt.fields.Sessions,
				FinishedGames: tt.fields.FinishedGames,
			}
			if err := s.DeleteSession(tt.args.sessionId); (err != nil) != tt.wantErr {
				t.Errorf("DeleteSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_EndGame(t *testing.T) {
	type fields struct {
		Router        *mux.Router
		Sessions      *sync.Map
		FinishedGames *cache.Cache
	}
	type args struct {
		sessionId string
		gameId    string
		playerId  string
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
			s := &Server{
				Router:        tt.fields.Router,
				Sessions:      tt.fields.Sessions,
				FinishedGames: tt.fields.FinishedGames,
			}
			if err := s.EndGame(tt.args.sessionId, tt.args.gameId, tt.args.playerId); (err != nil) != tt.wantErr {
				t.Errorf("EndGame() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_GetGameState(t *testing.T) {
	type fields struct {
		Router        *mux.Router
		Sessions      *sync.Map
		FinishedGames *cache.Cache
	}
	type args struct {
		sessionId string
		gameId    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Router:        tt.fields.Router,
				Sessions:      tt.fields.Sessions,
				FinishedGames: tt.fields.FinishedGames,
			}
			got, err := s.GetGameState(tt.args.sessionId, tt.args.gameId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGameState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGameState() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetSessionInfo(t *testing.T) {
	type fields struct {
		Router        *mux.Router
		Sessions      *sync.Map
		FinishedGames *cache.Cache
	}
	type args struct {
		sessionId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Router:        tt.fields.Router,
				Sessions:      tt.fields.Sessions,
				FinishedGames: tt.fields.FinishedGames,
			}
			got, got1, err := s.GetSessionInfo(tt.args.sessionId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSessionInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSessionInfo() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetSessionInfo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestServer_JoinGame(t *testing.T) {
	type fields struct {
		Router        *mux.Router
		Sessions      *sync.Map
		FinishedGames *cache.Cache
	}
	type args struct {
		sessionId  string
		gameId     string
		playerName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Router:        tt.fields.Router,
				Sessions:      tt.fields.Sessions,
				FinishedGames: tt.fields.FinishedGames,
			}
			got, err := s.JoinGame(tt.args.sessionId, tt.args.gameId, tt.args.playerName)
			if (err != nil) != tt.wantErr {
				t.Errorf("JoinGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JoinGame() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_ListOpenGames(t *testing.T) {
	type fields struct {
		Router        *mux.Router
		Sessions      *sync.Map
		FinishedGames *cache.Cache
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Router:        tt.fields.Router,
				Sessions:      tt.fields.Sessions,
				FinishedGames: tt.fields.FinishedGames,
			}
			if got := s.ListOpenGames(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListOpenGames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_NewSession(t *testing.T) {
	type fields struct {
		Router        *mux.Router
		Sessions      *sync.Map
		FinishedGames *cache.Cache
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Router:        tt.fields.Router,
				Sessions:      tt.fields.Sessions,
				FinishedGames: tt.fields.FinishedGames,
			}
			got, err := s.NewSession()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewSession() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_PlayMove(t *testing.T) {
	type fields struct {
		Router        *mux.Router
		Sessions      *sync.Map
		FinishedGames *cache.Cache
	}
	type args struct {
		sessionId string
		gameId    string
		playerId  string
		row       int
		col       int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Router:        tt.fields.Router,
				Sessions:      tt.fields.Sessions,
				FinishedGames: tt.fields.FinishedGames,
			}
			got, err := s.PlayMove(tt.args.sessionId, tt.args.gameId, tt.args.playerId, tt.args.row, tt.args.col)
			if (err != nil) != tt.wantErr {
				t.Errorf("PlayMove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PlayMove() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_authenticateSessionId(t *testing.T) {
	type fields struct {
		Router        *mux.Router
		Sessions      *sync.Map
		FinishedGames *cache.Cache
	}
	type args struct {
		sessionId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Session
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Router:        tt.fields.Router,
				Sessions:      tt.fields.Sessions,
				FinishedGames: tt.fields.FinishedGames,
			}
			got, err := s.authenticateSessionId(tt.args.sessionId)
			if (err != nil) != tt.wantErr {
				t.Errorf("authenticateSessionId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authenticateSessionId() got = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		// TODO: Add test cases.
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateGameInSession() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_EndGame(t *testing.T) {
	type fields struct {
		Id         string
		ActiveGame *game.Game
	}
	type args struct {
		gameId   string
		playerId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *game.Game
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Id:         tt.fields.Id,
				ActiveGame: tt.fields.ActiveGame,
			}
			got, err := s.EndGame(tt.args.gameId, tt.args.playerId)
			if (err != nil) != tt.wantErr {
				t.Errorf("EndGame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EndGame() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_GetGameState(t *testing.T) {
	type fields struct {
		Id         string
		ActiveGame *game.Game
	}
	type args struct {
		gameId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Id:         tt.fields.Id,
				ActiveGame: tt.fields.ActiveGame,
			}
			got, err := s.GetGameState(tt.args.gameId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGameState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGameState() got = %v, want %v", got, tt.want)
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
		name   string
		fields fields
		want   string
		want1  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Id:         tt.fields.Id,
				ActiveGame: tt.fields.ActiveGame,
			}
			got, got1 := s.GetSessionInfo()
			if got != tt.want {
				t.Errorf("GetSessionInfo() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetSessionInfo() got1 = %v, want %v", got1, tt.want1)
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
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if got != tt.want {
				t.Errorf("JoinGame() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_PlayMove(t *testing.T) {
	type fields struct {
		Id         string
		ActiveGame *game.Game
	}
	type args struct {
		gameId   string
		playerId string
		row      int
		col      int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *game.Game
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Id:         tt.fields.Id,
				ActiveGame: tt.fields.ActiveGame,
			}
			got, err := s.PlayMove(tt.args.gameId, tt.args.playerId, tt.args.row, tt.args.col)
			if (err != nil) != tt.wantErr {
				t.Errorf("PlayMove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlayMove() got = %v, want %v", got, tt.want)
			}
		})
	}
}
