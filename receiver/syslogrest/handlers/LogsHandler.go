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

const NUM_POST_KEY = "num"
const MSG_POST_KEY = "msg"
const APP_POST_KEY = "app"

const DEFAULT_NUM_STRING = "500"

type LogsHandler struct {
	store storage.LogStore
}

func NewLogsHandler(storage storage.LogStore) *LogsHandler {
	return &LogsHandler{
		store: storage,
	}
}

func (logsHandler *LogsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	start := time.Now()
	numRaw := req.PostFormValue(NUM_POST_KEY)
	msgRaw := req.PostFormValue(MSG_POST_KEY)
	appRaw := req.PostFormValue(APP_POST_KEY)

	if numRaw == "" {
		logging.Info.Println("Got request without num, using 500 as default")
		numRaw = DEFAULT_NUM_STRING
	}
	num, err := strconv.Atoi(numRaw)

	if err != nil {
		http.Error(
			res,
			fmt.Sprintf("Error parsing %s into an Integer. Error is: %s", numRaw, err.Error()),
			http.StatusBadRequest,
		)
		logging.Error.Printf("Error Handling Request: %s, Error was: %s", req, err.Error())
		return
	}

	var results []*parsing.Message
	if msgRaw != "" || appRaw != "" {
		logging.Debug.Printf("Returning %d entries that match: app: '%s', msg: '%s'", num, appRaw, msgRaw)
		results, err = logsHandler.store.Filter(appRaw, msgRaw, num)
	} else {
		logging.Debug.Printf("Nothing to search for, returning %d last entries", num)
		results, err = logsHandler.store.GetLatest(num)
	}

	if err != nil {
		http.Error(
			res,
			fmt.Sprintf("Couldn't retrieve Things, Sorry", err.Error()),
			http.StatusInternalServerError)
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
