package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/mrWinston/sysloco/receiver/logging"
	"github.com/mrWinston/sysloco/receiver/settings"
	"github.com/mrWinston/sysloco/receiver/storage"
	"github.com/mrWinston/sysloco/receiver/syslog"
	"github.com/mrWinston/sysloco/receiver/syslogrest"
)

var shutdown = make(chan bool, 1)
var osSig = make(chan os.Signal, 1)

func handleErr(err error) {
	if err != nil {
		logging.Error.Printf("%s", err)
		logging.Error.Printf("Shutting Down...", err)
		shutdown <- true
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
	logOpts := logging.DefaultOptions()
	options, err := settings.Parse()
	handleErr(err)
	// init Logging
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
		logging.Error.Fatal("Badger Store not implemented")
		//storageOpts := storage.BadgerDefaultOptions()
		//storageOpts.Dir = options.DbLocation
		//storageOpts.ValueDir = options.DbLocation
		//logStore, err = storage.NewBadgerStore(storageOpts)
		//handleErr(err)
	} else {
		logStore, err = storage.NewMemStore(options.DbLocation, 100)
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
	defer syslogServer.Stop()

	syslogRestServer, err := syslogrest.NewServer(options.HttpPort, options.HttpAddress, logStore)
	handleErr(err)
	err = syslogRestServer.Start()
	defer syslogRestServer.Stop()
	defer logStore.Release()

	// Wait for Shutdown Signal
	<-shutdown
	logging.Info.Println("Shutting Down...")
}
