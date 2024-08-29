package api

// request body and response body for APIs

type CreateNewSessionResp struct {
	SessionId string `json:"sessionId"`
}

type GetCurrentSessionResp struct {
	SessionId string `json:"sessionId"`
	GameId    string `json:"gameId"`
}

type CreateNewGameReq struct {
	PlayerName string `json:"playerName"`
}

type CreateNewGameResp struct {
	GameId   string `json:"gameId"`
	PlayerId string `json:"playerId"`
}

type ListOpenGamesResp struct {
	SessionIdAndGameIds map[string]string `json:"sessionIdAndGameIds"`
}

type GetGameStateResp struct {
	State string `json:"state"`
}

type JoinGameReq struct {
	PlayerName string `json:"playerName"`
}

type JoinGameResp struct {
	PlayerId string `json:"playerId"`
}

type EndGameReq struct {
	PlayerId string `json:"playerId"`
}

type PlayMoveReq struct {
	PlayerId string `json:"playerId"`
	Row      int    `json:"row"`
	Column   int    `json:"column"`
}

type PlayMoveResp struct {
	State string `json:"state"`
}
