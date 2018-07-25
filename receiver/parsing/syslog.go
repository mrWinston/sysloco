package parsing

import (
	"log"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Priv      string
	Version   int
	Timestamp time.Time
	Hostname  string
	Appname   string
	Procid    string
	Msgid     string
	Msg       string
}

func GetMsg(rawMessage string) *Message {
	msgarr := strings.Split(rawMessage, " ")
	privverarr := strings.SplitAfter(msgarr[0], ">")
	ver, err := strconv.Atoi(privverarr[1])
	if err != nil {
		log.Fatal("The Input didn't contain a proper version! input Was: ", rawMessage)
	}
	timestamp, _ := time.Parse(time.RFC3339, msgarr[1])

	return &Message{
		Priv:      privverarr[0],
		Version:   ver,
		Timestamp: timestamp,
		Hostname:  msgarr[2],
		Appname:   msgarr[3],
		Procid:    msgarr[4],
		Msgid:     msgarr[5],
		Msg:       strings.TrimSpace(strings.Join(msgarr[7:], " ")),
	}

}
