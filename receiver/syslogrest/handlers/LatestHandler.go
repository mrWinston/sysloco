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

var numPostKey = "num"

type LatestHandler struct {
	store storage.LogStore
}

func NewLatestHandler(storage storage.LogStore) *LatestHandler {
	return &LatestHandler{
		store: storage,
	}
}

func (latestHandler *LatestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	start := time.Now()
	numString := req.PostFormValue(numPostKey)

	if numString == "" {
		http.Error(res, fmt.Sprintf("\"%s\" parameter not provided!", numPostKey), http.StatusBadRequest)
		return
	}
	num, err := strconv.Atoi(numString)
	if err != nil {
		http.Error(res, fmt.Sprintf("Error parsing %s into an Integer. Error is: %s", numString, err.Error()), http.StatusBadRequest)
		logging.Error.Printf("Error Handling Request: %s, Error was: %s", req, err.Error())
		return
	}
	results, err := latestHandler.store.GetLatest(num)
	if err != nil {
		http.Error(res, fmt.Sprintf("Couldn't retrieve Things, Sorry", err.Error()), http.StatusInternalServerError)
		logging.Error.Printf("Error Handling Request: %s, Error was: %s", req, err.Error())
		return
	}
	logging.Debug.Printf("Found %d Results when Getting the Last.", len(results))
	var wrapStruct = struct {
		Result []*parsing.Message
	}{
		Result: results,
	}

	jsReturn, err := json.Marshal(&wrapStruct)
	if err != nil {
		http.Error(res, fmt.Sprintf("Oopsie"), http.StatusInternalServerError)
		logging.Error.Printf("Error handling request: %s Error was: %s", req, err.Error())
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	res.Write(jsReturn)
	logging.Debug.Printf("Answer Time: %s", time.Since(start))
}
