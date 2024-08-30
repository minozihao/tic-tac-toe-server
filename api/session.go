package api

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/minozihao/tic-tac-toe-server/game"
)

// global in memory variables

var (
	DefaultSize = 1000
)

// predefined errors

var (
	GameIdNotMatchErr        = errors.New("game id not match")
	SessionIdAuthErr         = errors.New("authentication error. invalid session id")
	NoActiveGameInSessionErr = errors.New("no active game in session")
)

type Session struct {
	Id         string
	ActiveGame *game.Game
}

// CreateGameInSession create an active game in the session and returns the game object
func (s *Session) CreateGameInSession(playerName string) (*game.Game, error) {
	gameFactory := game.NewGameFactory{}
	newGame := gameFactory.CreateGame(playerName)
	s.ActiveGame = newGame
	return newGame, nil
}

func (s *Session) GetSessionInfo() (string, string) {
	if s.ActiveGame == nil {
		return s.Id, ""
	}
	return s.Id, s.ActiveGame.Id
}

func (s *Session) GetGameState(gameId string) (string, error) {
	if s.ActiveGame == nil {
		return "", NoActiveGameInSessionErr
	}
	if gameId != s.ActiveGame.Id {
		return "", GameIdNotMatchErr
	}
	return s.ActiveGame.ShowGameState(s.Id), nil
}

// JoinGame join the game and returns a player2 id
func (s *Session) JoinGame(gameId, playerName string) (string, error) {
	if s.ActiveGame == nil {
		return "", NoActiveGameInSessionErr
	}
	player2Id := uuid.NewString()
	err := s.ActiveGame.Join(gameId, player2Id, playerName)
	if err != nil {
		return "", err
	}
	return player2Id, nil
}

// EndGame remove the game from active game field in session and return the game pointer
func (s *Session) EndGame(gameId, playerId string) (*game.Game, error) {
	if s.ActiveGame == nil {
		return nil, NoActiveGameInSessionErr
	}
	if err := s.ActiveGame.EndGame(gameId, playerId); err != nil {
		return nil, err
	}
	g := s.ActiveGame
	s.ActiveGame = nil
	return g, nil
}

// PlayMove play a legal move and returns the game pointer
func (s *Session) PlayMove(gameId string, playerId string, row int, col int) (*game.Game, error) {
	if s.ActiveGame == nil {
		return nil, NoActiveGameInSessionErr
	}
	if s.ActiveGame == nil {
		return nil, NoActiveGameInSessionErr
	}
	if s.ActiveGame.Id != gameId {
		return nil, GameIdNotMatchErr
	}
	if err := s.ActiveGame.Move(playerId, row, col); err != nil {
		return nil, err
	}
	return s.ActiveGame, nil
}

// Functions for controller to call

// NewSession create a new session and register into in memory sync map sessions
func (s *Server) NewSession() (string, error) {
	sid := uuid.NewString()
	var ss = Session{
		Id:         sid,
		ActiveGame: nil,
	}
	// check sessions size. limit 1000
	var count int
	s.Sessions.Range(func(k, v any) bool {
		count++
		return true
	})
	if count >= DefaultSize {
		return "", errors.New("sessions limit 1000 reached. please wait for new space")
	}
	s.Sessions.Store(sid, &ss)
	return sid, nil
}

// DeleteSession delete the session from InMemSession syncMap if id match
func (s *Server) DeleteSession(sessionId string) error {
	_, err := s.authenticateSessionId(sessionId)
	if err != nil {
		return err
	}
	s.Sessions.Delete(sessionId)
	return nil
}

// GetSessionInfo returns current session id and active/open game id
func (s *Server) GetSessionInfo(sessionId string) (string, string, error) {
	session, err := s.authenticateSessionId(sessionId)
	if err != nil {
		return "", "", err
	}
	if session.ActiveGame == nil {
		return session.Id, "", nil
	}
	return session.Id, session.ActiveGame.Id, nil
}

// CreateGame create an open game in session, returns game id and player 1 id for the host
func (s *Server) CreateGame(sessionId string, playerName string) (string, string, error) {
	session, err := s.authenticateSessionId(sessionId)
	if err != nil {
		return "", "", err
	}
	ga, err := session.CreateGameInSession(playerName)
	if err != nil {
		return "", "", err
	}
	return ga.Id, ga.Player1Id, nil
}

// ListOpenGames returns a map of sessionId vs open game id for all sessions
func (s *Server) ListOpenGames() map[string]string {
	var openGames = make(map[string]string)
	s.Sessions.Range(func(sessionId, session any) bool {
		s := session.(*Session)
		if s.ActiveGame != nil && s.ActiveGame.Player2Id == "" {
			openGames[sessionId.(string)] = s.ActiveGame.Id
		}
		return true
	})
	return openGames
}

// GetGameState show game state of finished game or active game for the given session id and game id
func (s *Server) GetGameState(sessionId, gameId string) (string, error) {
	// check finished game
	cacheKey := fmt.Sprintf("%s_%s", sessionId, gameId)
	g, found := s.FinishedGames.Get(cacheKey)
	if found {
		ga := g.(*game.Game)
		return ga.ShowGameState(sessionId), nil
	}
	// check active games in session
	session, err := s.authenticateSessionId(sessionId)
	if err != nil {
		return "", err
	}

	output, err := session.GetGameState(gameId)
	if err != nil {
		return "", err
	}
	return output, nil
}

// JoinGame join a game in a session returns the id for player 2
func (s *Server) JoinGame(sessionId, gameId, playerName string) (string, error) {
	session, err := s.authenticateSessionId(sessionId)
	if err != nil {
		return "", err
	}
	playerId, err := session.JoinGame(gameId, playerName)
	if err != nil {
		return "", err
	}
	return playerId, nil
}

// EndGame change state of the game, remove the game from session and add it to the finished game cache
func (s *Server) EndGame(sessionId, gameId, playerId string) error {
	session, err := s.authenticateSessionId(sessionId)
	if err != nil {
		return err
	}
	g, err := session.EndGame(gameId, playerId)
	if err != nil {
		return err
	}
	cacheKey := fmt.Sprintf("%s_%s", sessionId, gameId)
	s.FinishedGames.Set(cacheKey, g, 0)
	return nil
}

// PlayMove play a legal move and returns the game state
func (s *Server) PlayMove(sessionId string, gameId string, playerId string, row int, col int) (string, error) {
	session, err := s.authenticateSessionId(sessionId)
	if err != nil {
		return "", err
	}
	g, err := session.PlayMove(gameId, playerId, row, col)
	if err != nil {
		return "", err
	}
	// if game finished, we need to remove the game from session and add it to finishedGame cache
	if g.State.End {
		session.ActiveGame = nil
		if s.FinishedGames.ItemCount() > DefaultSize {
			s.FinishedGames.DeleteExpired()
			if s.FinishedGames.ItemCount() > DefaultSize {
				log.Print("warning.finishedGames cache max size reached")
			}
		}
		cacheKey := fmt.Sprintf("%s_%s", sessionId, gameId)
		s.FinishedGames.Set(cacheKey, g, 0)
	}
	return g.ShowGameState(sessionId), nil
}

func (s *Server) authenticateSessionId(sessionId string) (*Session, error) {
	ss, ok := s.Sessions.Load(sessionId)
	if !ok {
		return nil, SessionIdAuthErr
	}
	session := ss.(*Session)
	if session == nil {
		return nil, SessionIdAuthErr
	}
	return session, nil
}
