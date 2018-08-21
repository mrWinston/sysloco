package syslogrest

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/mrWinston/sysloco/receiver/logging"
	"github.com/mrWinston/sysloco/receiver/storage"
	"github.com/mrWinston/sysloco/receiver/syslogrest/handlers"
)

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
	rootHandler := handlers.NewRootHandler(syslogRestServer.store)
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
