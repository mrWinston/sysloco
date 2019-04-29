package syslogrest

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mrWinston/sysloco/receiver/logging"
	"github.com/mrWinston/sysloco/receiver/storage"
	"github.com/mrWinston/sysloco/receiver/syslogrest/handlers"
)

type rootHandler struct {
	filterHandler *handlers.FilterHandler
	latestHandler *handlers.LatestHandler
	logsHandler   *handlers.LogsHandler
}

func newRootHandler(storage storage.LogStore) *rootHandler {
	return &rootHandler{
		filterHandler: handlers.NewFilterHandler(storage),
		latestHandler: handlers.NewLatestHandler(storage),
		logsHandler:   handlers.NewLogsHandler(storage),
	}
}

func (rootHandler *rootHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	logging.Error.Printf("Receiving Request for %s, Path is: %s", req.URL, req.URL.Path)
	path := strings.Split(req.URL.Path, "/")[1:]
	switch path[0] {
	case "filter":
		rootHandler.filterHandler.ServeHTTP(res, req)
		return
	case "latest":
		rootHandler.latestHandler.ServeHTTP(res, req)
		return
	case "logs":
		rootHandler.logsHandler.ServeHTTP(res, req)
		return
	default:
		http.Error(res, fmt.Sprintf("404 Not Found"), http.StatusNotFound)
		logging.Error.Printf("Path not Found: %s", path[0])
		return
	}

}

// SyslogRestServer is an instance of the REST API for the Receiver. You can
// use this Struct directly, but it is advised to use syslogrest.NewServer()
// instead, as it provides input validation
type SyslogRestServer struct {
	port       int
	address    string
	store      storage.LogStore
	httpServer *http.Server
	running    bool
}

func NewServer(port int, address string, store storage.LogStore) (*SyslogRestServer, error) {
	return &SyslogRestServer{
		port:    port,
		address: address,
		store:   store,
		running: false,
	}, nil
}

func (syslogRestServer *SyslogRestServer) Start() error {
	if syslogRestServer.running {
		return errors.New("Server already started")
	}

	logging.Debug.Printf("Starting server with Addr: %s, Port: %d", syslogRestServer.address, syslogRestServer.port)
	rootHandler := newRootHandler(syslogRestServer.store)
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", syslogRestServer.address, syslogRestServer.port),
		Handler:      rootHandler,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		ErrorLog:     logging.Error,
	}
	logging.Debug.Printf("Created a Server")
	logging.Debug.Printf("%s:%d", syslogRestServer.address, syslogRestServer.port)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", syslogRestServer.address, syslogRestServer.port))
	logging.Debug.Printf("Created the Listener")
	if err != nil {
		return err
	}
	logging.Debug.Printf("Starting to serve")
	go httpServer.Serve(listener)
	logging.Debug.Printf("Now Serving")
	syslogRestServer.httpServer = httpServer
	syslogRestServer.running = true
	return nil
}

func (syslogRestServer *SyslogRestServer) Stop() error {
	if syslogRestServer.httpServer == nil {
		return errors.New("Server isn't Started!")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := syslogRestServer.httpServer.Shutdown(ctx)
	syslogRestServer.running = false
	return err
}
