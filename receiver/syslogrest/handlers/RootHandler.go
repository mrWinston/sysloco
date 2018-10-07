package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mrWinston/sysloco/receiver/logging"
	"github.com/mrWinston/sysloco/receiver/storage"
)

type RootHandler struct {
	filterHandler *FilterHandler
	latestHandler *LatestHandler
}

func NewRootHandler(storage storage.LogStore) *RootHandler {
	return &RootHandler{
		filterHandler: NewFilterHandler(storage),
		latestHandler: NewLatestHandler(storage),
	}
}

func (rootHandler *RootHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	logging.Error.Printf("Receiving Request for %s, Path is: %s", req.URL, req.URL.Path)
	path := strings.Split(req.URL.Path, "/")[1:]
	switch path[0] {
	case "filter":
		rootHandler.filterHandler.ServeHTTP(res, req)
		return
	case "latest":
		rootHandler.latestHandler.ServeHTTP(res, req)
	default:
		http.Error(res, fmt.Sprintf("404 Not Found"), http.StatusNotFound)
		logging.Error.Printf("Path not Found: %s", path[0])
		return
	}

}
