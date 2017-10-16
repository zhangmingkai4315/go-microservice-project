package service

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/unrolled/render"
	engine "github.com/zhangmingkai4315/go-microservice-project/engine"
)

func createMatchHandler(formatter *render.Render) http.HandlerFunc{
	return func(w http.ResponseWriter,req *http.Request){
		payload,_:=ioutil.ReadAll(req.Body)
		fmt.Printf("%+v",string(payload))
		var newMatchRequest newMatchRequest
		err:=json.Unmarshal(payload,&newMatchRequest)
		fmt.Printf("%+v",newMatchRequest)
		if err != nil||(newMatchRequest.isValid()==false){
			formatter.JSON(w,http.StatusBadRequest,"Fail to parse match request")
			return
		}
		newMatch := engine.NewMatch(
						newMatchRequest.GridSize,
						newMatchRequest.PlayerBlack,
						newMatchRequest.PlayerWhite)
		w.Header().Add("Location","/matches/"+newMatch.ID)
		formatter.JSON(w,http.StatusCreated,struct{Test string}{"hello"})
	}
}

