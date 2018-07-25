package syslogrest

import (
	"errors"
	"net/http"

	"github.com/dgraph-io/badger"
)

//TODO: Remove the Opts struct
type Opts struct {
	Port    int
	Address string
}

type Server struct {
	Opts Opts
	DB   *badger.DB
}

func showLast100(w *http.ResponseWriter, r *http.Request) {

}
func Init(opts Opts) (*Server, error) {
	if opts.Port > 65535 || opts.Port <= 0 {
		return nil, errors.New("The http-Port is invalid")
	}
	return nil, nil
}
