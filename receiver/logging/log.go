package logging

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

const defaultLogMask = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile

//TODO: Remove The Opts struct, and just pass the values
type Opts struct {
	Level     int
	DebugPath string
	InfoPath  string
	ErrorPath string
}

func (o Opts) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{\n")
	s := reflect.ValueOf(&o).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		buffer.WriteString(fmt.Sprintf("\t %s = %v,\n", typeOfT.Field(i).Name, f.Interface()))
	}
	buffer.WriteString("}\n")
	return buffer.String()
}

var (
	Debug   *log.Logger = log.New(os.Stdout, "[DEBUG] ", defaultLogMask)
	Info    *log.Logger = log.New(os.Stdout, "[INFO] ", defaultLogMask)
	Error   *log.Logger = log.New(os.Stdout, "[ERROR] ", defaultLogMask)
	LogOpts Opts
)

func Init(opts Opts) error {
	if opts.Level > 2 || opts.Level < 0 {
		return errors.New("Level needs to be between 0 and 2 ( Including both )")
	}
	paths := []string{opts.ErrorPath, opts.InfoPath, opts.DebugPath}
	prefs := []string{"[ERROR] ", "[INFO]  ", "[DEBUG] "}
	loggers := make([]*log.Logger, 3)

	for i := 0; i < 3; i++ {
		if paths[i] != "" {
			f, err := os.OpenFile(paths[i], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				return err
			}
			if i <= opts.Level {
				loggers[i] = log.New(io.MultiWriter(f, os.Stdout), prefs[i], defaultLogMask)
			} else {
				loggers[i] = log.New(f, prefs[i], defaultLogMask)
			}
		} else {
			if i <= opts.Level {
				loggers[i] = log.New(os.Stdout, prefs[i], defaultLogMask)
			} else {
				loggers[i] = log.New(ioutil.Discard, prefs[i], defaultLogMask)
			}
		}
	}
	Debug = loggers[2]
	Info = loggers[1]
	Error = loggers[0]

	// Debug out for stuff
	Debug.Println("Initialized Logging")
	Debug.Println("Using the following Options:")
	Debug.Print(opts)

	return nil
}

func DefaultOptions() *Opts {
	return &Opts{
		Level:     1,
		DebugPath: "",
		InfoPath:  "",
		ErrorPath: "",
	}
}
