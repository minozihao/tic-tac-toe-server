package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

type Server struct {
	*mux.Router
	// store session as value and session id as key
	Sessions *sync.Map
	// FinishedGames a cache with 1 min default expiration and purges expired items every 2 mins
	// store finsihed games with key being sessionId + '_' + gameId
	FinishedGames *cache.Cache
}

func NewServer() *Server {
	s := &Server{
		Router:        mux.NewRouter(),
		Sessions:      &sync.Map{},
		FinishedGames: cache.New(1*time.Minute, 2*time.Minute),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.HandleFunc("/session", s.createNewSession()).Methods("POST")
	s.HandleFunc("/session", s.getCurrentSession()).Methods("GET")
	s.HandleFunc("/session", s.endSession()).Methods("DELETE")

	// game handlers
	s.HandleFunc("/games", s.createNewGame()).Methods("POST")
	s.HandleFunc("/games", s.listOpenGames()).Methods("GET")
	s.HandleFunc("/games/{gameId}", s.getGameState()).Methods("GET")
	s.HandleFunc("/games/{gameId}/join", s.joinGame()).Methods("POST")
	s.HandleFunc("/games/{gameId}/play", s.playMove()).Methods("POST")
	s.HandleFunc("/games/{gameId}", s.endGame()).Methods("DELETE")
}

// createNewSession create a new session (should return a token/ID which can be used for authentication)
func (s *Server) createNewSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		sid, err := s.NewSession()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var resp = &CreateNewSessionResp{
			SessionId: sid,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

// getCurrentSession get current session ID and active/open game ID
func (s *Server) getCurrentSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.Header.Get("Authorization")
		if sessionId == "" {
			http.Error(w, errors.New("no sessionId found in header authorization").Error(), http.StatusUnauthorized)
			return
		}
		sessionId, gameId, err := s.GetSessionInfo(sessionId)
		if errors.Is(err, SessionIdAuthErr) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var resp = &GetCurrentSessionResp{
			SessionId: sessionId,
			GameId:    gameId,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

// endSession delete session
func (s *Server) endSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.Header.Get("Authorization")
		if sessionId == "" {
			http.Error(w, errors.New("no sessionId found in header authorization").Error(), http.StatusUnauthorized)
			return
		}
		err := s.DeleteSession(sessionId)
		if errors.Is(err, SessionIdAuthErr) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

// createNewGame create a new game (sets the current game ID for authenticated session)
func (s *Server) createNewGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.Header.Get("Authorization")
		if sessionId == "" {
			http.Error(w, errors.New("no sessionId found in header authorization").Error(), http.StatusUnauthorized)
			return
		}
		var body CreateNewGameReq
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		if err := d.Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		gameId, playerId, err := s.CreateGame(sessionId, body.PlayerName)
		if errors.Is(err, SessionIdAuthErr) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var resp = &CreateNewGameResp{
			GameId:   gameId,
			PlayerId: playerId,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

// listOpenGames list open games in all sessions
func (s *Server) listOpenGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		openGames := s.ListOpenGames()
		var resp = &ListOpenGamesResp{
			SessionIdAndGameIds: openGames,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

// getGameState get the game state
func (s *Server) getGameState() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.Header.Get("Authorization")
		if sessionId == "" {
			http.Error(w, errors.New("no sessionId found in header authorization").Error(), http.StatusUnauthorized)
			return
		}
		gameId, found := mux.Vars(r)["gameId"]
		if !found {
			http.Error(w, errors.New("game id not found in path").Error(), http.StatusBadRequest)
			return
		}

		state, err := s.GetGameState(sessionId, gameId)
		if errors.Is(err, SessionIdAuthErr) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var resp = &GetGameStateResp{
			State: state,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

// joinGame join an open game (sets the current game ID for the authenticated session)
func (s *Server) joinGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.Header.Get("Authorization")
		if sessionId == "" {
			http.Error(w, errors.New("no sessionId found in header authorization").Error(), http.StatusUnauthorized)
			return
		}
		gameId, found := mux.Vars(r)["gameId"]
		if !found {
			http.Error(w, errors.New("game id not found in path").Error(), http.StatusBadRequest)
			return
		}

		var body JoinGameReq
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		if err := d.Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		playerId, err := s.JoinGame(sessionId, gameId, body.PlayerName)
		if errors.Is(err, SessionIdAuthErr) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var resp = &JoinGameResp{
			PlayerId: playerId,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

// playMove play a legal move
func (s *Server) playMove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.Header.Get("Authorization")
		if sessionId == "" {
			http.Error(w, errors.New("no sessionId found in header authorization").Error(), http.StatusUnauthorized)
			return
		}
		gameId, found := mux.Vars(r)["gameId"]
		if !found {
			http.Error(w, errors.New("game id not found in path").Error(), http.StatusBadRequest)
			return
		}

		var body PlayMoveReq
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		if err := d.Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		state, err := s.PlayMove(sessionId, gameId, body.PlayerId, body.Row, body.Column)
		if errors.Is(err, SessionIdAuthErr) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var resp = &PlayMoveResp{
			State: state,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}

// endGame end game
func (s *Server) endGame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.Header.Get("Authorization")
		if sessionId == "" {
			http.Error(w, errors.New("no sessionId found in header authorization").Error(), http.StatusUnauthorized)
			return
		}
		gameId, found := mux.Vars(r)["gameId"]
		if !found {
			http.Error(w, errors.New("game id not found in path").Error(), http.StatusBadRequest)
			return
		}

		var body EndGameReq
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()
		if err := d.Decode(&body); err != nil {
			http.Error(w, fmt.Errorf("invalid request body, %w", err).Error(), http.StatusBadRequest)
			return
		}

		err := s.EndGame(sessionId, gameId, body.PlayerId)
		if errors.Is(err, SessionIdAuthErr) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}
