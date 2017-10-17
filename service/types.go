package service

import (
	engine "github.com/zhangmingkai4315/go-microservice-project/engine"
	"gopkg.in/mgo.v2/bson"
	"net/url"
)

type newMatchRequest struct {
	GridSize    int    `json:"gridsize"`
	PlayerWhite string `json:"playerWhite"`
	PlayerBlack string `json:"playerBlack"`
}

func (request newMatchRequest) isValid() (valid bool) {
	valid = true
	if request.GridSize != 19 && request.GridSize != 13 && request.GridSize != 9 {
		valid = false
	}
	if request.PlayerBlack == "" || request.PlayerWhite == "" {
		valid = false
	}
	return valid
}

type newMatchResponse struct {
	ID          string   `json:"id"`
	StartedAt   int64    `json:"started_at"`
	GridSize    int      `json:"gridsize"`
	PlayerWhite string   `json:"playerWhite"`
	PlayerBlack string   `json:"playerBlack"`
	Turn        int      `json:"turn,omitempty"`
	GameBoard   [][]byte `json:"gameboard"`
}

func (m *newMatchResponse) copyMatch(match engine.Match) {
	m.ID = match.ID
	m.StartedAt = match.StartTime.Unix()
	m.GridSize = match.GridSize
	m.PlayerWhite = match.PlayerWhite
	m.PlayerBlack = match.PlayerBlack
	m.Turn = match.TurnCount
	m.GameBoard = match.GameBoard.Positions
}

type matchDetailsResponse struct {
	ID          string   `json:"id"`
	StartedAt   int64    `json:"started_at"`
	GridSize    int      `json:"gridsize"`
	PlayerWhite string   `json:"playerWhite"`
	PlayerBlack string   `json:"playerBlack"`
	Turn        int      `json:"turn,omitempty"`
	GameBoard   [][]byte `json:"gameboard"`
}

func (m *matchDetailsResponse) copyMatch(match engine.Match) {
	m.ID = match.ID
	m.StartedAt = match.StartTime.Unix()
	m.GridSize = match.GridSize
	m.PlayerWhite = match.PlayerWhite
	m.PlayerBlack = match.PlayerBlack
	m.Turn = match.TurnCount
	m.GameBoard = match.GameBoard.Positions
}

type matchRepository interface {
	addMatch(match engine.Match) (err error)
	getMatches() (matches []engine.Match, err error)
	getMatch(id string) (match engine.Match, err error)
}

type RequestParams struct {
	RawQuery url.Values `json:"raw_query"`
	Q        bson.M     `json:"selector"`
	S        bson.M     `json:"scope"`
	L        int        `json:"limit"`
	F        int        `json:"offset"`
}

type dbConfig struct {
	HostURI          string
	DBCollectionName string
}
