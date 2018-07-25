package storage

import (
	"os"
	"testing"
	"time"

	"github.com/mrWinston/sysloco/receiver/parsing"
	"github.com/stretchr/testify/assert"
)

var sleeptime = time.Second / 100
var memStoreFile = "/tmp/memstoretest"
var msgs = []*parsing.Message{
	&parsing.Message{
		Priv:      "2",
		Version:   1,
		Timestamp: time.Now(),
		Hostname:  "mrGizmo",
		Appname:   "Superapp",
		Procid:    "IDIDIDI",
		Msgid:     "hlasfhlisfhli",
		Msg:       "Device: /dev/sda [SAT], SMART Usage Attribute: 194 Temperature_Celsius changed from 47 to 48",
	},
	&parsing.Message{
		Priv:      "2",
		Version:   1,
		Timestamp: time.Now(),
		Hostname:  "mrGizmo",
		Appname:   "Superapp",
		Procid:    "IDIDIDI",
		Msgid:     "hlasfhlisfhli",
		Msg:       "kdeconnect.core: TCP connection done (i'm the existing device)",
	},
	&parsing.Message{
		Priv:      "2",
		Version:   1,
		Timestamp: time.Now(),
		Hostname:  "mrGizmo",
		Appname:   "Superapp",
		Procid:    "IDIDIDI",
		Msgid:     "hlasfhlisfhli",
		Msg:       `127.0.0.1 - - [28/Nov/2017:16:10:52 +0100] "GET /en/all-countries-ajax/ HTTP/1.1" 200 9585 "http://my.3yd/en/psconfig/appearance/" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.75 Safari/537.36"`,
	},
	&parsing.Message{
		Priv:      "2",
		Version:   1,
		Timestamp: time.Now(),
		Hostname:  "mrGizmo",
		Appname:   "Superapp",
		Procid:    "IDIDIDI",
		Msgid:     "hlasfhlisfhli",
		Msg: `2017/11/28 15:51:59 [error] 26940#26940: *8 open() "/tmp/button3d/assets/images/favicon.png" failed (2: No such file or directory), client: 127.0.0.1, server: , request: "GET /assets/images/favicon.png HTTP/1.1", host: "my.3yd", referrer: "http://my.3yd/en/dashboard/"
`,
	},
	&parsing.Message{
		Priv:      "2",
		Version:   1,
		Timestamp: time.Now(),
		Hostname:  "mrGizmo",
		Appname:   "Superapp",
		Procid:    "IDIDIDI",
		Msgid:     "hlasfhlisfhli",
		Msg:       `localhost - - [22/Jul/2018:14:37:19 +0200] "POST / HTTP/1.1" 200 149 Cancel-Subscription successful-ok`,
	},
}

func checkEquality(first *parsing.Message, second *parsing.Message, t *testing.T) bool {
	return false
}

func initMemStore(remove bool, t *testing.T) *MemStore {
	if remove {
		os.Remove(memStoreFile)
	}

	memStore, err := NewMemStore(memStoreFile)
	assert.Nil(t, err)

	fileInfo, err := os.Stat(memStoreFile)

	assert.Nil(t, err)
	assert.Equal(t, fileInfo.IsDir(), false)
	return memStore

}

func TestNewMemStore(t *testing.T) {
	// test, if using a nonexisting file works
	initMemStore(true, t)
	// test, if using the existing file works
	initMemStore(false, t)
}

func TestStore(t *testing.T) {
	memStore := initMemStore(true, t)

	for i := 0; i < len(msgs); i++ {
		memStore.Store(*msgs[i])
		time.Sleep(sleeptime)
		assert.Equal(t, msgs[i].Msg, memStore.store[i].Msg)
	}

}

func TestGetLatest(t *testing.T) {
	memStore := initMemStore(true, t)
	memStore.Store(*msgs[0])
	time.Sleep(sleeptime)
	res, err := memStore.GetLatest(1)
	assert.Nil(t, err)
	assert.Equal(t, res[0].Msg, msgs[0].Msg)

	for i := 0; i < len(msgs); i++ {
		memStore.Store(*msgs[i])
		time.Sleep(sleeptime)
	}
	res, err = memStore.GetLatest(len(msgs))

	for i := 0; i < len(msgs); i++ {
		assert.Equal(t, res[i].Msg, msgs[len(msgs)-(i+1)].Msg)
	}
}

func TestFilter(t *testing.T) {
	memStore := initMemStore(true, t)
	for i := 0; i < len(msgs); i++ {
		memStore.Store(*msgs[i])
		time.Sleep(sleeptime)
	}
	res, err := memStore.Filter("my.3yd")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))

	assert.Equal(t, res[0].Msg, msgs[3].Msg)
	assert.Equal(t, res[1].Msg, msgs[2].Msg)
}
