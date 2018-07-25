package settings

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"reflect"
)

type Opts struct {
	LogLevel       int
	DebugPath      string
	InfoPath       string
	ErrorPath      string
	HttpPort       int
	HttpAddress    string
	SyslogPort     int
	SyslogAddress  string
	SyslogProtocol string
	SyslogFormat   string
	DbEngine       string
	DbLocation     string
	Config         string
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

type opt struct {
	name     string
	empty    interface{}
	argtype  string
	helpmsg  string
	required bool
}

var optdesc = []opt{
	{
		name:     "v",
		empty:    -1,
		argtype:  "int",
		helpmsg:  "Verbosity: [0-2]",
		required: false,
	},
	{
		name:     "logfile-debug",
		empty:    "",
		argtype:  "string",
		helpmsg:  "Path to Debug Log File( Leave empty for no file logging )",
		required: false,
	},
	{
		name:     "logfile-info",
		empty:    "",
		argtype:  "string",
		helpmsg:  "Path to Info Log File( Leave empty for no file logging )",
		required: false,
	},
	{
		name:     "logfile-err",
		empty:    "",
		argtype:  "string",
		helpmsg:  "Path to Error Log File( Leave empty for no file logging )",
		required: false,
	},
	{
		name:     "http-port",
		empty:    -1,
		argtype:  "int",
		helpmsg:  "Port of Http-Rest server",
		required: false,
	},
	{
		name:     "http-address",
		empty:    "",
		argtype:  "string",
		helpmsg:  "Ip, on which to listen for http-Rest-Requests",
		required: false,
	},
	{
		name:     "syslog-port",
		empty:    10001,
		argtype:  "int",
		helpmsg:  "Port of the Syslog Server [10001]",
		required: false,
	},
	{
		name:     "syslog-address",
		empty:    "0.0.0.0",
		argtype:  "string",
		helpmsg:  "The Address, on which the syslog server listens [0.0.0.0]",
		required: false,
	},
	{
		name:     "syslog-proto",
		empty:    "",
		argtype:  "string",
		helpmsg:  "The Transport protocol for Syslog: UDP or TCP [ UDP ]",
		required: false,
	},
	{
		name:     "syslog-format",
		empty:    "",
		argtype:  "string",
		helpmsg:  "The Syslog Protocol version to use [rfc5424]",
		required: false,
	},
	{
		name:     "db-engine",
		empty:    "memory",
		argtype:  "string",
		helpmsg:  "The DB Engine: badger or memory",
		required: false,
	},
	{
		name:     "db-loc",
		empty:    "",
		argtype:  "string",
		helpmsg:  "The Dir of the badger db Files or the Persistency file for In Mem DB",
		required: true,
	},
	{
		name:     "conf",
		empty:    "",
		argtype:  "string",
		helpmsg:  "Path to a config file",
		required: false,
	},
}

func DefaultOptions() Opts {
	return Opts{
		LogLevel:       1,
		DebugPath:      "",
		InfoPath:       "",
		ErrorPath:      "",
		HttpPort:       80,
		HttpAddress:    "0.0.0.0",
		SyslogPort:     10001,
		SyslogAddress:  "0.0.0.0",
		SyslogProtocol: "UDP",
		SyslogFormat:   "rfc5424",
		DbEngine:       "memory",
		DbLocation:     "",
		Config:         "",
	}
}

func Parse() (*Opts, error) {
	cmdlineargs := make([]interface{}, len(optdesc), len(optdesc))
	for i, option := range optdesc {
		if empty, ok := option.empty.(string); ok {
			cmdlineargs[i] = flag.String(option.name, empty, option.helpmsg)
		} else if empty, ok := option.empty.(int); ok {
			cmdlineargs[i] = flag.Int(option.name, empty, option.helpmsg)
		} else if empty, ok := option.empty.(bool); ok {
			cmdlineargs[i] = flag.Bool(option.name, empty, option.helpmsg)
		} else {
			return nil, errors.New("Types didn't match when parsing: " + option.name)
		}
	}
	flag.Parse()
	// parsing conf file
	options := DefaultOptions()
	if confpath, ok := cmdlineargs[len(cmdlineargs)-1].(*string); ok && *confpath != "" {
		raw, err := ioutil.ReadFile(*confpath)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(raw, &options)
	}
	// huge block for settings the correct values

	s := reflect.ValueOf(&options).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.Type().String() == "int" {
			if val, ok := cmdlineargs[i].(*int); ok && *val != -1 {
				s.Field(i).SetInt(int64(*val))
			}
			if optdesc[i].required && f.Interface() == -1 {
				return nil, errors.New(fmt.Sprint(typeOfT.Field(i).Name, " Needs to be set!"))
			}
		} else if f.Type().String() == "string" {
			if val, ok := cmdlineargs[i].(*string); ok && *val != "" {
				s.Field(i).SetString(*val)
			}
			if optdesc[i].required && f.Interface() == "" {
				return nil, errors.New(fmt.Sprint(typeOfT.Field(i).Name, " Needs to be set!"))
			}
		}
	}

	return &options, nil
}
