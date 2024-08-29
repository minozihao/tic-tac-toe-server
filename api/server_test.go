package api

import (
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
