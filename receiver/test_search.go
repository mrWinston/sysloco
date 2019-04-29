package main

import (
	"bufio"
	"encoding/json"
	"os"
	"regexp"
	"time"

	"github.com/mrWinston/sysloco/receiver/logging"
	"github.com/mrWinston/sysloco/receiver/parsing"
)

func main_disabled() {
	logging.Debug.Println("Yo")

	start := time.Now()
	appExp, err := regexp.Compile("ui")
	persistencyFile := "./test/large"
	logging.Info.Println("Loading Existing Store File from ", persistencyFile)
	file, err := os.Open(persistencyFile)
	if err != nil {
		logging.Error.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	logging.Debug.Printf("Took %s to initialize the reader", time.Since(start))

	for scanner.Scan() {
		var msg parsing.Message
		jsonMsg := scanner.Text()
		err := json.Unmarshal([]byte(jsonMsg), &msg)
		if err != nil {
			logging.Info.Println("The Store is corrupted, couldn't Unmarshal line: ", err)
			continue
		}
		appExp.MatchString(msg.Msg)
	}

	logging.Debug.Printf("Took %s for everything", time.Since(start))

}
