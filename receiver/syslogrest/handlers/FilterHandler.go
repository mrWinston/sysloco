package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mrWinston/sysloco/receiver/logging"
	"github.com/mrWinston/sysloco/receiver/parsing"
	"github.com/mrWinston/sysloco/receiver/storage"
)

var regexPostKey = "filter"

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
	regexString := req.PostFormValue(regexPostKey)
	if regexString == "" {
		http.Error(res, fmt.Sprintf("\"%s\" parameter not provided!", regexPostKey), http.StatusBadRequest)
		return
	}

	results, err := filterHandler.store.Filter(regexString)
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
