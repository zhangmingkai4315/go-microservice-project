package service

import (
	"strings"
	"io/ioutil"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unrolled/render"
)
var (
	formatter = render.New(render.Options{
		IndentJSON:true,
	})
)

const (
	fakeMatchLocationResult = "/matches/5a003b78-409e-4452-b456-a6f0dcee05bd"
)

func TestCreateMatchWithCorrectObject(t *testing.T){
	client:=&http.Client{}
	server:=httptest.NewServer(http.HandlerFunc(createMatchHandler(formatter)))
	defer server.Close()

	body := []byte("{\n  \"gridsize\": 19,\"playerWhite\":\"Mike\",\"playerBlack\":\"Alice\"}")

	req,err:=http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
	if err != nil{
		t.Errorf("Error in creating POST request for createMatchHandler: %v",err)
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err!=nil{
		t.Errorf("Error in send POST request for createMatchHandler: %v",err)
	}

	defer res.Body.Close()

	payload,err:=ioutil.ReadAll(res.Body)
	if err!=nil{
		t.Errorf("Error in read POST response for createMatchHandler:%v",err)
	}	

	if res.StatusCode!=http.StatusCreated{
		t.Errorf("Expected 201 status code. Got:%d",res.StatusCode)
	}

	loc, headerOk:=res.Header["Location"]
	if !headerOk{
		t.Errorf("Location header is not set!")
	}else{
		if !strings.Contains(loc[0],"/matches/"){
			t.Errorf("Location Header Should Contain '/matches/'")
		}
		if len(loc[0])!=len(fakeMatchLocationResult){
			t.Errorf("Location length is not good for parse")
		}
	}
	t.Logf("Payload:%s", string(payload))
}

func TestCreateMatchHandlerWithBadRequest(t *testing.T){
	client:=&http.Client{}
	server:=httptest.NewServer(http.HandlerFunc(createMatchHandler(formatter)))
	defer server.Close()

	body1 := []byte("{\n  \"test\":\"testing for notvalid json object!\"}")
	req,err:=http.NewRequest("POST", server.URL, bytes.NewBuffer(body1))
	if err != nil{
		t.Errorf("Error in creating POST request for createMatchHandler: %v",err)
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err!=nil{
		t.Errorf("Error in send POST request for createMatchHandler: %v",err)
	}

	defer res.Body.Close()
	payload,err:=ioutil.ReadAll(res.Body)
	if err!=nil{
		t.Errorf("Error in read POST response for createMatchHandler:%v",err)
	}

	if res.StatusCode != http.StatusBadRequest{
		t.Errorf("Sending invalid JSON should result in a bad request from server:Got %d",res.StatusCode)
	}
	t.Logf("Payload:%s", string(payload))
	body2 := []byte("test")
	req,err=http.NewRequest("POST", server.URL, bytes.NewBuffer(body2))
	if err != nil{
		t.Errorf("Error in creating POST request for createMatchHandler: %v",err)
	}

	req.Header.Add("Content-Type", "application/json")
	res, err = client.Do(req)
	if err!=nil{
		t.Errorf("Error in send POST request for createMatchHandler: %v",err)
	}
	
	defer res.Body.Close()
	payload,err=ioutil.ReadAll(res.Body)
	if err!=nil{
	t.Errorf("Error in read POST response for createMatchHandler: %v",err)
	}
	if res.StatusCode != http.StatusBadRequest{
		t.Errorf("Sending invalid JSON should result in a bad request from server:Got %d",res.StatusCode)
	}
	t.Logf("Payload:%s", string(payload))
}

