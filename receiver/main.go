package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrWinston/sysloco/receiver/logging"
	"github.com/mrWinston/sysloco/receiver/settings"
	"github.com/mrWinston/sysloco/receiver/storage"
	"github.com/mrWinston/sysloco/receiver/syslog"
)

var shutdown = make(chan bool, 1)
var osSig = make(chan os.Signal, 1)

func handleErr(err error) {
	if err != nil {
		if logging.Error != nil {
			logging.Error.Fatal(err)
		} else {
			log.Fatal(err)
		}
	}
}

func handleOsSignal() {
	sig := <-osSig
	logging.Info.Println("Received ", sig)
	logging.Info.Println("Shutting Down")
	shutdown <- true
}

func main() {
	// init Cmdline Args
	signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM)
	go handleOsSignal()
	options, err := settings.Parse()
	handleErr(err)

	// init Logging
	logOpts := logging.DefaultOptions()
	logOpts.Level = options.LogLevel
	logOpts.DebugPath = options.DebugPath
	logOpts.InfoPath = options.InfoPath
	logOpts.ErrorPath = options.ErrorPath
	err = logging.Init(*logOpts)
	handleErr(err)

	logging.Debug.Println("Initializing with these Options:")
	logging.Debug.Println("\n", options)

	var logStore storage.LogStore
	if options.DbEngine == "badger" {
		// Settings up Badger Storage
		storageOpts := storage.BadgerDefaultOptions()
		storageOpts.Dir = options.DbLocation
		storageOpts.ValueDir = options.DbLocation
		logStore, err = storage.NewBadgerStore(storageOpts)
		handleErr(err)
	} else {
		logStore, err = storage.NewMemStore(options.DbLocation)
		handleErr(err)
	}

	// Setup Syslog Server Listener
	syslogOptions := syslog.DefaultOpts()
	syslogOptions.Port = options.SyslogPort
	syslogOptions.Ip = options.SyslogAddress
	syslogServer, err := syslog.New(syslogOptions)
	handleErr(err)
	syslogServer.DB = logStore
	go syslogServer.Start()

	// Wait for Shutdown Signal
	<-shutdown
	logging.Info.Println("Shutting Down...")
	syslogServer.Stop()
	logStore.Release()

}
