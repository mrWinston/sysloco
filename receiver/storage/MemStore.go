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

// The MemStore Struct provides Methods for interacting with the Memory Backed
// Log Storage
type MemStore struct {
	storeFile string
	store     []*parsing.Message
	msgChan   chan *parsing.Message
	stop      bool
}

// NewMemStore Returns a new Instance of of the MemStore struct. It either
// creates a new persistencyFile or loads all log msg from an existing one
func NewMemStore(persistencyFile string) (*MemStore, error) {
	// file doens't exist
	start := time.Now()
	logging.Debug.Println("Loading MemStore Persistency File ")
	var file *os.File
	if _, err := os.Stat(persistencyFile); os.IsNotExist(err) {
		logging.Debug.Println("Creating new File at: ", persistencyFile)
		file, err = os.Create(persistencyFile)
		if err != nil {
			return nil, err
		}
	} else {
		logging.Debug.Println("Loading Existing File from ", persistencyFile)
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
		storeFile: persistencyFile,
		store:     lines,
		msgChan:   make(chan *parsing.Message),
		stop:      false,
	}
	go memStore.listenForMsg()

	logging.Debug.Printf("Took me %s to load the file", time.Since(start))
	return memStore, nil
}

func (memStore *MemStore) listenForMsg() {
	for !memStore.stop {
		msg := <-memStore.msgChan
		memStore.store = append(memStore.store, msg)
	}
}

// Store puts the parsing.Message into the memStore. Adding happens
// asyncronously, so calling GetLatest right after adding will likely not
// return the just added msg
func (memStore *MemStore) Store(msg parsing.Message) error {
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
func (memStore *MemStore) Filter(filter string) ([]*parsing.Message, error) {
	var ret []*parsing.Message
	matchExp, err := regexp.Compile(filter)
	if err != nil {
		return nil, err
	}
	storeSize := len(memStore.store)
	for i := storeSize - 1; i >= 0; i-- {
		curMsg := memStore.store[i]
		if matchExp.MatchString(curMsg.Msg) {
			ret = append(ret, curMsg)
		}
	}
	return ret, nil
}
