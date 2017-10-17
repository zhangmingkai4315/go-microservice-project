package service

import (
	"bytes"
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	engine "github.com/zhangmingkai4315/go-microservice-project/engine"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

const (
	fakeMatchLocationResult = "/matches/5a003b78-409e-4452-b456-a6f0dcee05bd"
)

func TestCreateMatchWithCorrectObject(t *testing.T) {
	client := &http.Client{}
	repo := NewInMemoryRepository()
	server := httptest.NewServer(http.HandlerFunc(createMatchHandler(formatter, repo)))
	defer server.Close()

	body := []byte("{\n  \"gridsize\": 19,\"playerWhite\":\"Mike\",\"playerBlack\":\"Alice\"}")

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in send POST request for createMatchHandler: %v", err)
	}

	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error in read POST response for createMatchHandler:%v", err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected 201 status code. Got:%d", res.StatusCode)
	}

	loc, headerOk := res.Header["Location"]
	if !headerOk {
		t.Errorf("Location header is not set!")
	} else {
		if !strings.Contains(loc[0], "/matches/") {
			t.Errorf("Location Header Should Contain '/matches/'")
		}
		if len(loc[0]) != len(fakeMatchLocationResult) {
			t.Errorf("Location length is not good for parse")
		}
	}
}

func TestCreateMatchHandlerWithBadRequest(t *testing.T) {
	client := &http.Client{}
	repo := NewInMemoryRepository()
	server := httptest.NewServer(http.HandlerFunc(createMatchHandler(formatter, repo)))
	defer server.Close()

	body1 := []byte("{\n  \"test\":\"testing for notvalid json object!\"}")
	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body1))
	if err != nil {
		t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in send POST request for createMatchHandler: %v", err)
	}

	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error in read POST response for createMatchHandler:%v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Sending invalid JSON should result in a bad request from server:Got %d", res.StatusCode)
	}
	body2 := []byte("test")
	req, err = http.NewRequest("POST", server.URL, bytes.NewBuffer(body2))
	if err != nil {
		t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	res, err = client.Do(req)
	if err != nil {
		t.Errorf("Error in send POST request for createMatchHandler: %v", err)
	}

	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error in read POST response for createMatchHandler: %v", err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Sending invalid JSON should result in a bad request from server:Got %d", res.StatusCode)
	}

}

func TestGetMatchListReturnsEmptyArrayForNoMatches(t *testing.T) {
	client := &http.Client{}
	repo := NewInMemoryRepository()
	server := httptest.NewServer(http.HandlerFunc(getMatchListHandler(formatter, repo)))
	defer server.Close()
	req, _ := http.NewRequest("GET", server.URL, nil)

	resp, err := client.Do(req)

	if err != nil {
		t.Error("Errored when sending request to the server", err)
		return
	}

	defer resp.Body.Close()
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("Failed to read response from server", err)
	}

	var matchList []newMatchResponse
	err = json.Unmarshal(payload, &matchList)
	if err != nil {
		t.Errorf("Could not unmarshal payload into []newMatchResponse slice")
	}

	if len(matchList) != 0 {
		t.Errorf("Expected an empty list of match responses, got %d", len(matchList))
	}
}

func TestGetMatchListReturnsWhatsInRepository(t *testing.T) {
	client := &http.Client{}
	repo := NewInMemoryRepository()

	repo.addMatch(engine.NewMatch(19, "black", "white"))
	repo.addMatch(engine.NewMatch(13, "bl", "wh"))
	repo.addMatch(engine.NewMatch(19, "b", "w"))
	server := httptest.NewServer(http.HandlerFunc(getMatchListHandler(formatter, repo)))
	defer server.Close()
	req, _ := http.NewRequest("GET", server.URL, nil)

	resp, err := client.Do(req)

	if err != nil {
		t.Error("Errored when sending request to the server", err)
		return
	}

	defer resp.Body.Close()
	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("Failed to read response from server", err)
	}

	var matchList []newMatchResponse
	err = json.Unmarshal(payload, &matchList)
	if err != nil {
		t.Errorf("Could not unmarshal payload into []newMatchResponse slice")
	}

	repoMatches, err := repo.getMatches()
	if err != nil {
		t.Errorf("Unexpected error in getMatches(): %s", err)
	}
	if len(matchList) != len(repoMatches) {
		t.Errorf("Match response size should have equaled repo size, sizes were: %d and %d", len(matchList), len(repoMatches))
	}

	for idx := 0; idx < 3; idx++ {
		if matchList[idx].GridSize != repoMatches[idx].GridSize {
			t.Errorf("Gridsize mismatch at index %d. Got %d and %d", idx, matchList[idx].GridSize, repoMatches[idx].GridSize)
		}
		if matchList[idx].PlayerBlack != matchList[idx].PlayerBlack {
			t.Errorf("PlayerBlack mismatch at index %d. Got %s and %s", idx, matchList[idx].PlayerBlack, repoMatches[idx].PlayerBlack)
		}
		if matchList[idx].PlayerWhite != matchList[idx].PlayerWhite {
			t.Errorf("PlayerWhite mismatch at index %d. Got %s and %s", idx, matchList[idx].PlayerWhite, repoMatches[idx].PlayerWhite)
		}
	}
}

func MakeTestServer(repo matchRepository) *negroni.Negroni {
	server := negroni.New()
	mx := mux.NewRouter()
	initRoutes(mx, formatter, repo)
	server.UseHandler(mx)
	return server
}
func TestGetMatchDetailsReturns404ForNonexistentMatch(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	repo := NewInMemoryRepository()
	server := MakeTestServer(repo)
	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/matches/1234", nil)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected %v; received %v", http.StatusNotFound, recorder.Code)
	}
}

func TestGetMatchDetailsReturnsExistingMatch(t *testing.T) {
	var (
		request  *http.Request
		recorder *httptest.ResponseRecorder
	)

	repo := NewInMemoryRepository()
	server := MakeTestServer(repo)

	targetMatch := engine.NewMatch(19, "black", "white")
	repo.addMatch(targetMatch)
	targetMatchID := targetMatch.ID
	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest("GET", "/matches/"+targetMatchID, nil)
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected %v; received %v", http.StatusOK, recorder.Code)
	}

	var match matchDetailsResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &match)
	if err != nil {
		t.Errorf("Error unmarshaling match details: %s", err)
	}
	if match.GridSize != 19 {
		t.Errorf("Expected match gridsize to be 19; received %d", match.GridSize)
	}
}
