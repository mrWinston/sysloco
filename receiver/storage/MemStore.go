package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/mrWinston/sysloco/receiver/logging"
	"github.com/mrWinston/sysloco/receiver/parsing"
)

type WrappingMessageArray struct {
	size  int
	index int
	store []*parsing.msg
}

func NewWrappingMessageArray(size int) *WrappingMessageArray {
	if size <= 0 {
		return &WrappingMessageArray

}

// The MemStore Struct provides Methods for interacting with the Memory Backed
// Log Storage
type MemStore struct {
	maximumLines   int
	maximumAgeDays int
	cleanUpTick    *time.Ticker
	storeFile      string
	store          []*parsing.Message
	msgChan        chan *parsing.Message
	stop           chan bool
}

// NewMemStore Returns a new Instance of of the MemStore struct. It either
// creates a new persistencyFile or loads all log msg from an existing one
func NewMemStore(persistencyFile string, maximumLines int, maximumAgeDays int) (*MemStore, error) {
	// file doens't exist
	start := time.Now()
	var file *os.File
	if _, err := os.Stat(persistencyFile); os.IsNotExist(err) {
		logging.Info.Println("Creating new Store at: ", persistencyFile)
		file, err = os.Create(persistencyFile)
		if err != nil {
			return nil, err
		}
	} else {
		logging.Info.Println("Loading Existing Store File from ", persistencyFile)
		file, err = os.Open(persistencyFile)
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []*parsing.Message

	for scanner.Scan() {
		var msg parsing.Message
		jsonMsg := scanner.Text()
		err := json.Unmarshal([]byte(jsonMsg), &msg)
		if err != nil {
			logging.Info.Println("The Store is corrupted, couldn't Unmarshal line: ", err)
			continue
		}
		lines = append(lines, &msg)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	var memStore = &MemStore{
		maximumLines:   maximumLines,
		maximumAgeDays: maximumAgeDays,
		cleanUpTick:    time.NewTicker(5 * time.Second),
		storeFile:      persistencyFile,
		store:          lines,
		msgChan:        make(chan *parsing.Message),
		stop:           make(chan bool),
	}
	go memStore.listenForMsg()
	go memStore.runOldMessageCleaner()

	logging.Debug.Printf("Took me %s to load the file", time.Since(start))
	return memStore, nil
}

func (memStore *MemStore) runOldMessageCleaner() {
	for {
		select {
		case <-memStore.stop:
			logging.Debug.Println("Stopping Cleanup Routine")
			return
		case <-memStore.cleanUpTick.C:
			logging.Debug.Println("Cleaning...")
		}
	}
}

func (memStore *MemStore) listenForMsg() {
	for {
		select {
		case <-memStore.stop:
			logging.Debug.Println("Stopped listening for messages")
			return
		case msg := <-memStore.msgChan:
			memStore.store = append(memStore.store, msg)
		}
	}
}

// Store puts the parsing.Message into the memStore. Adding happens
// asyncronously, so calling GetLatest right after adding will likely not
// return the just added msg
func (memStore *MemStore) Store(msg parsing.Message) error {
	logging.Debug.Printf("Receiving Msg")
	memStore.msgChan <- &msg
	return nil
}

// GetLatest Returns the Last added elements to the store, In Reverse order ( from
// newest to oldest )
func (memStore *MemStore) GetLatest(number int) ([]*parsing.Message, error) {
	length := len(memStore.store)
	var results []*parsing.Message
	var maxNum int
	if number < length {
		maxNum = number
	} else {
		maxNum = length
	}

	for i := 0; i < maxNum; i++ {
		msg := memStore.store[length-1-i]
		results = append(results, msg)
	}

	return results, nil
}

// Release Stops the Listen Loop for Messages and write the Storage back to the file
func (memStore *MemStore) Release() error {
	memStore.stop <- true

	memStore.cleanUpTick.Stop()
	close(memStore.stop)
	file, err := os.Create(memStore.storeFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, msg := range memStore.store {
		byteLine, err := json.Marshal(msg)
		if err != nil {
			logging.Info.Println("Couldn's Transform msg to json: ", err)
			continue
		}
		fmt.Fprintln(writer, string(byteLine))
	}
	return writer.Flush()
}

// Filter Returns all Elements that match the provided regex in reverse order ( from
// newest to oldest )
func (memStore *MemStore) Filter(appFilter string, msgFilter string, num int) ([]*parsing.Message, error) {
	logging.Debug.Printf("Filtering for %s and %s with %s results ...", appFilter, msgFilter, num)
	var ret []*parsing.Message

	appExp, err := regexp.Compile(appFilter)
	if err != nil {
		logging.Info.Printf("Received an Invalid AppFilter: \"%s\"", appFilter)
		return nil, err
	}
	msgExp, err := regexp.Compile(msgFilter)
	if err != nil {
		logging.Info.Printf("Received an Invalid MsgFilter: \"%s\"", msgFilter)
		return nil, err
	}

	storeSize := len(memStore.store)
	for i := storeSize - 1; i >= 0 && len(ret) <= num; i-- {
		curMsg := memStore.store[i]
		if !msgExp.MatchString(curMsg.Msg) {
			continue
		} else if appExp.MatchString(curMsg.Appname) {
			ret = append(ret, curMsg)
		}
	}
	logging.Debug.Printf("Found %d Results", len(ret))
	return ret, nil
}
