package service

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	engine "github.com/zhangmingkai4315/go-microservice-project/engine"
	"io/ioutil"
	"net/http"
)

func createMatchHandler(formatter *render.Render, repo matchRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var newMatchRequest newMatchRequest
		err := json.Unmarshal(payload, &newMatchRequest)
		if err != nil || (newMatchRequest.isValid() == false) {
			formatter.JSON(w, http.StatusBadRequest, "Fail to parse match request")
			return
		}
		newMatch := engine.NewMatch(
			newMatchRequest.GridSize,
			newMatchRequest.PlayerBlack,
			newMatchRequest.PlayerWhite)
		repo.addMatch(newMatch)
		w.Header().Add("Location", "/matches/"+newMatch.ID)
		formatter.JSON(w, http.StatusCreated, struct{ Test string }{"hello"})
	}
}

func getMatchListHandler(formatter *render.Render, repo matchRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		repoMatches, err := repo.getMatches()
		if err != nil {
			formatter.JSON(w, http.StatusNotFound, err.Error())
		} else {
			matches := make([]newMatchResponse, len(repoMatches))
			for idx, match := range repoMatches {
				matches[idx].copyMatch(match)
			}
			formatter.JSON(w, http.StatusOK, matches)
		}
	}
}

func getMatchDetailsHandler(formatter *render.Render, repo matchRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		matchId := vars["id"]
		match, err := repo.getMatch(matchId)
		if err != nil {
			formatter.JSON(w, http.StatusNotFound, err.Error())
		} else {
			var mdr matchDetailsResponse
			mdr.copyMatch(match)
			formatter.JSON(w, http.StatusOK, &mdr)
		}
	}
}
