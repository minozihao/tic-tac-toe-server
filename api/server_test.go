package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateNewSession(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/session", nil)
	w := httptest.NewRecorder()
	NewServer().createNewSession()(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(string(data)) < 15 || !strings.Contains(string(data), "{\"sessionId\":") {
		t.Error("Expected a generated uuid string session id")
	}
}

func TestGetCurrentSession_AuthError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/session", nil)
	req.Header.Set("Authorization", "abc")
	w := httptest.NewRecorder()
	NewServer().getCurrentSession()(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Println(string(data))
	if string(data) != "authentication error. invalid session id\n" {
		t.Error("unexpected error message in response")
	}
}

func TestGetCurrentSession_Success(t *testing.T) {
	// create a session
	var s = NewServer()
	req := httptest.NewRequest(http.MethodPost, "/session", nil)
	w := httptest.NewRecorder()
	s.createNewSession()(w, req)
	res := w.Result()
	defer res.Body.Close()
	var resp CreateNewSessionResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sessionId := resp.SessionId
	expected := fmt.Sprintf("{\"sessionId\":\"%s\",\"gameId\":\"\"}\n", sessionId)
	// call get current session with sessionId
	req2 := httptest.NewRequest(http.MethodGet, "/session", nil)
	req2.Header.Set("Authorization", sessionId)
	w2 := httptest.NewRecorder()
	s.getCurrentSession()(w2, req2)
	res2 := w2.Result()
	defer res2.Body.Close()
	data, err := io.ReadAll(res2.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if (string(data)) != expected {
		t.Errorf("result not match. actual: %s. expect: %s", string(data), expected)
	}

}

func TestCreateNewGame_AuthErr(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/games", nil)
	req.Header.Set("Authorization", "abc")
	w := httptest.NewRecorder()
	NewServer().getCurrentSession()(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if string(data) != "authentication error. invalid session id\n" {
		t.Error("unexpected error message in response")
	}
}

func TestCreateNewGame_Success(t *testing.T) {
	// create a session
	var s = NewServer()
	req := httptest.NewRequest(http.MethodPost, "/session", nil)
	w := httptest.NewRecorder()
	s.createNewSession()(w, req)
	res := w.Result()
	defer res.Body.Close()
	var resp CreateNewSessionResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sessionId := resp.SessionId

	// create game in the session
	var body = CreateNewGameReq{
		PlayerName: "bob",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		t.Errorf("unexpect error %s", err.Error())
	}
	req2 := httptest.NewRequest(http.MethodPost, "/games", &buf)
	req2.Header.Set("Authorization", sessionId)
	w2 := httptest.NewRecorder()
	s.createNewGame()(w2, req2)
	res2 := w2.Result()
	defer res2.Body.Close()
	var newResp CreateNewGameResp
	if err := json.NewDecoder(res2.Body).Decode(&newResp); err != nil {
		t.Errorf("unexpect error %s", err.Error())
	}
	if newResp.GameId == "" || newResp.PlayerId == "" {
		t.Error("expect generated uuid for game id and player id")
	}
}

func TestListOpenGames(t *testing.T) {
	// create a session
	var s = NewServer()
	req := httptest.NewRequest(http.MethodPost, "/session", nil)
	w := httptest.NewRecorder()
	s.createNewSession()(w, req)
	res := w.Result()
	defer res.Body.Close()
	var resp CreateNewSessionResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sessionId := resp.SessionId

	// create game in the session
	var body = CreateNewGameReq{
		PlayerName: "bob",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		t.Errorf("unexpect error %s", err.Error())
	}
	req2 := httptest.NewRequest(http.MethodPost, "/games", &buf)
	req2.Header.Set("Authorization", sessionId)
	w2 := httptest.NewRecorder()
	s.createNewGame()(w2, req2)

	// check list open games
	req3 := httptest.NewRequest(http.MethodGet, "/games", &buf)
	req3.Header.Set("Authorization", sessionId)
	w3 := httptest.NewRecorder()
	s.listOpenGames()(w3, req3)
	res3 := w3.Result()
	defer res3.Body.Close()
	var newResp ListOpenGamesResp
	if err := json.NewDecoder(res3.Body).Decode(&newResp); err != nil {
		t.Errorf("unexpect error %s", err.Error())
	}
	if len(newResp.SessionIdAndGameIds) != 1 {
		t.Error("expect 1 open game in the session")
	}
}
