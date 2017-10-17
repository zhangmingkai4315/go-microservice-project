package service

import (
	"fmt"
	"github.com/cloudnativego/cfmgo"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
)

func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	// repo := initRepository(true, "")
	repo := initRepository(false, &dbConfig{HostURI: "mongodb://127.0.0.1:27017/gogo", DBCollectionName: "matches"})
	initRoutes(mx, formatter, repo)
	n.UseHandler(mx)
	return n
}
func initRepository(memoryType bool, db *dbConfig) (repo matchRepository) {
	if memoryType == true {
		return NewInMemoryRepository()
	} else {
		matchCollection := cfmgo.Connect(cfmgo.NewCollectionDialer, db.HostURI, db.DBCollectionName)
		fmt.Printf("Connecting to MongoDB service: %s...\n", db.HostURI)
		repo = NewMongoMatchRepository(matchCollection)
		return repo
	}
}
func initRoutes(mx *mux.Router, formatter *render.Render, repo matchRepository) {
	mx.HandleFunc("/test", testHandler(formatter)).Methods("GET")
	mx.HandleFunc("/matches", createMatchHandler(formatter, repo)).Methods("POST")
	mx.HandleFunc("/matches", getMatchListHandler(formatter, repo)).Methods("GET")
	mx.HandleFunc("/matches/{id}", getMatchDetailsHandler(formatter, repo)).Methods("GET")
}

func testHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{
			"This is a test func",
		})
	}
}
