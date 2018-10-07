package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mrWinston/sysloco/receiver/logging"
	"github.com/mrWinston/sysloco/receiver/parsing"
	"github.com/mrWinston/sysloco/receiver/storage"
)

var msgFilterKey = "msg"
var appFilterKey = "app"
var numKey = "num"

type FilterHandler struct {
	store storage.LogStore
}

func NewFilterHandler(storage storage.LogStore) *FilterHandler {
	return &FilterHandler{
		store: storage,
	}
}

func (filterHandler *FilterHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	start := time.Now()

	appFilterRegex := req.PostFormValue(appFilterKey)
	msgFilterRegex := req.PostFormValue(msgFilterKey)
	numString := req.PostFormValue(numKey)

	if numString == "" {
		numString = "200"
	}

	num, err := strconv.Atoi(numString)
	if err != nil {
		http.Error(res, fmt.Sprintf("Error parsing %s into an Integer. Error is: %s", numString, err.Error()), http.StatusBadRequest)
		logging.Error.Printf("Error Handling Request: %s, Error was: %s", req, err.Error())
		return
	}

	if msgFilterRegex == "" && appFilterRegex == "" {
		http.Error(res, fmt.Sprintf("No parameter not provided! Either \"%s\" or \"%s\" needs to be set!", appFilterKey, msgFilterKey), http.StatusBadRequest)
		return
	}

	results, err := filterHandler.store.Filter(appFilterRegex, msgFilterRegex, num)
	if err != nil {
		http.Error(res, fmt.Sprintf("Regex Invalid: %s", err.Error()), http.StatusBadRequest)
	}
	logging.Debug.Printf("Found %d Results when filtering.", len(results))
	var wrapStruct = struct {
		Result []*parsing.Message
	}{
		Result: results,
	}

	jsReturn, err := json.Marshal(&wrapStruct)
	if err != nil {
		http.Error(res, fmt.Sprintf("Oopsie"), http.StatusInternalServerError)
		logging.Error.Printf("Error handling request: %s Error was: %s", req, err)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	res.Write(jsReturn)
	logging.Debug.Printf("Answer Time: %s", time.Since(start))
}
