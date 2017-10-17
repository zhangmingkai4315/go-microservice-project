package service

import (
	"errors"
	"fmt"
	"github.com/cloudnativego/cfmgo"
	"github.com/cloudnativego/cfmgo/params"
	engine "github.com/zhangmingkai4315/go-microservice-project/engine"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type mongoMatchRepository struct {
	Collection cfmgo.Collection
}

type matchRecord struct {
	RecordID    bson.ObjectId `bson:"_id,omitempty" json:"id"`
	MatchID     string        `bson:"match_id" json:"match_id"`
	TurnCount   int           `bson:"turn_count" json:"turn_count"`
	GridSize    int           `bson:"grid_size" json:"grid_size"`
	StartTime   string        `bson:"start_time" json:"start_time"`
	GameBoard   [][]byte      `bson:"game_board" json:"game_board"`
	PlayerBlack string        `bson:"player_black" json:"player_black"`
	PlayerWhite string        `bson:"player_white" json:"player_white"`
}

func NewMongoMatchRepository(col cfmgo.Collection) (repo *mongoMatchRepository) {
	repo = &mongoMatchRepository{
		Collection: col,
	}
	return
}

func (r *mongoMatchRepository) addMatch(match engine.Match) (err error) {
	r.Collection.Wake()
	mr := convertMatchToMatchRecord(match)
	_, err = r.Collection.UpsertID(mr.RecordID, mr)
	return
}

func (r *mongoMatchRepository) getMatch(id string) (match engine.Match, err error) {
	r.Collection.Wake()
	theMatch, err := r.getMongoMatch(id)
	if err == nil {
		match = convertMatchRecordToMatch(theMatch)
	}
	return
}

func (r *mongoMatchRepository) getMongoMatch(id string) (mongoMatch matchRecord, err error) {
	var matches []matchRecord
	query := bson.M{"match_id": id}
	params := &params.RequestParams{
		Q: query,
	}

	count, err := r.Collection.Find(params, &matches)
	if count == 0 {
		err = errors.New("Match not found")
	}
	if err == nil {
		mongoMatch = matches[0]
	}
	return
}

func (r *mongoMatchRepository) getMatches() (matches []engine.Match, err error) {
	r.Collection.Wake()
	var mr []matchRecord
	_, err = r.Collection.Find(cfmgo.ParamsUnfiltered, &mr)
	if err == nil {
		matches = make([]engine.Match, len(mr))
		for k, v := range mr {
			matches[k] = convertMatchRecordToMatch(v)
		}
	}
	return
}

func convertMatchToMatchRecord(m engine.Match) (mr *matchRecord) {
	mr = &matchRecord{
		RecordID:    bson.NewObjectId(),
		MatchID:     m.ID,
		TurnCount:   m.TurnCount,
		GridSize:    m.GridSize,
		StartTime:   m.StartTime.Format("2006-01-02 15:04:05"),
		GameBoard:   m.GameBoard.Positions,
		PlayerBlack: m.PlayerBlack,
		PlayerWhite: m.PlayerWhite,
	}
	return
}

func convertMatchRecordToMatch(mr matchRecord) (m engine.Match) {
	t, err := time.Parse("2006-01-02 15:04:05", mr.StartTime)
	if err != nil {
		fmt.Printf("Error parsing time value in Match Record: %v", err)
	} else {
		m = engine.Match{
			ID:          mr.MatchID,
			TurnCount:   mr.TurnCount,
			GridSize:    mr.GridSize,
			StartTime:   t,
			GameBoard:   engine.GameBoard{Positions: mr.GameBoard},
			PlayerBlack: mr.PlayerBlack,
			PlayerWhite: mr.PlayerWhite,
		}
	}
	return
}
