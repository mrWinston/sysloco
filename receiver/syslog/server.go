package syslog

import (
	"errors"
	"fmt"
	"net"

	"github.com/mrWinston/sysloco/receiver/logging"
	"github.com/mrWinston/sysloco/receiver/parsing"
	"github.com/mrWinston/sysloco/receiver/storage"
)

// Opts holds the Options for the Syslogging server. You can set the following
// Values:
// 		* Port: The Default port on which the server listens
// 		* IP: The Ip address on which the server listens
// 		* BufferSize: The Internal Read buffer size. should be large enough to
// 			hold the biggest message ( In Bytes )
// 		* Protocol: The Transport Protocol to use ( 'UDP' or 'TCP' )
//TODO: Remove the Opts struct, and just pass the values to the New Function
type Opts struct {
	Port       int
	Ip         string
	BufferSize int
	Protocol   string
}

// The Server Struct holds the Methods for Starting the Syslogging server and
// the Options created. It provides Simple Methods for starting and Stopping
// the Server
type Server struct {
	Opts          Opts
	UdpConnection *net.UDPConn
	TcpConnection *net.TCPConn
	stopped       bool
	DB            storage.LogStore
	msgChan       chan ([]byte)
	stop          chan bool
}

const BUFFER_SIZE = 1000
const NUM_WORKERS = 10

// New creates a new Server instance from the given Opts.  Create the Opts
// beforehand with the DefaultOpts method.
func New(opts Opts) (*Server, error) {
	if opts.Port <= 0 || opts.Port >= 65535 {
		return nil, errors.New("Specify a Port between 1 and 65535 ( including )")
	}
	if opts.BufferSize <= 0 {
		return nil, errors.New("The BufferSize should be set to be greater than zero")
	}
	if opts.Protocol != "UDP" {
		return nil, errors.New("Only UDP is currently supported")
	}

	ServerAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", opts.Ip, opts.Port))
	if err != nil {
		return nil, err
	}
	ServerCon, err := net.ListenUDP("udp", ServerAddress)
	if err != nil {
		return nil, err
	}

	server := &Server{
		Opts:          opts,
		stopped:       false,
		UdpConnection: ServerCon,
		msgChan:       make(chan []byte, BUFFER_SIZE),
		stop:          make(chan bool),
	}

	logging.Debug.Printf("Starting %d Workers for syslog intake...", NUM_WORKERS)
	for i := 0; i < NUM_WORKERS; i++ {
		go server.handleMessages()
	}

	return server, nil
}

func DefaultOpts() Opts {
	return Opts{
		Port:       10001,
		Ip:         "0.0.0.0",
		BufferSize: 20480,
		Protocol:   "UDP",
	}
}

func (s *Server) Start() {

	for !s.stopped {
		buf := make([]byte, s.Opts.BufferSize)
		n, _, err := s.UdpConnection.ReadFromUDP(buf)
		if err != nil {
			logging.Info.Println("Got an Error while receiving message: ")
			logging.Info.Println(err)
		} else {
			logging.Info.Printf("got msg, waiting are: %v\n", len(s.msgChan))
			s.msgChan <- buf[0:n]
			//			go s.handleMessage(buf[0:n])
		}
	}
}

func (s *Server) Stop() {
	s.stopped = true
}

func (s *Server) handleMessages() {
	for {
		select {
		case <-s.stop:
			logging.Debug.Println("Stop Listening for syslog msgs")
			return
		case raw := <-s.msgChan:
			msg := parsing.GetMsg(string(raw))
			err := s.DB.Store(*msg)
			if err != nil {
				logging.Error.Println("GOT AN ERROR:", err)
			}
		}
	}
}

//func (s *Server) handleMessage(raw []byte) {
//	msg := parsing.GetMsg(string(raw))
//	err := s.DB.Store(*msg)
//	if err != nil {
//		logging.Error.Println("GOT AN ERROR:", err)
//	}
//}
