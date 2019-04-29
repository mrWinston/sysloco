package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/mrWinston/sysloco/receiver/logging"
	"github.com/mrWinston/sysloco/receiver/parsing"
)

// The percentage of fill ratio that should be achieved after cleaning
const AFTER_CLEAN_RATIO = 0.8

// The MemStore Struct provides Methods for interacting with the Memory Backed
// Log Storage
type MemStore struct {
	cleanUpTick      *time.Ticker
	maximumLines     int
	msgChan          chan *parsing.Message
	stop             chan bool
	store            []*parsing.Message
	storeFile        string
	cleanUpWaitGroup sync.WaitGroup
	storeMutex       *sync.RWMutex
}

// NewMemStore Returns a new Instance of of the MemStore struct. It either
// creates a new persistencyFile or loads all log msg from an existing one
func NewMemStore(persistencyFile string, maximumLines int) (*MemStore, error) {
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
		maximumLines: maximumLines,
		cleanUpTick:  time.NewTicker(5 * time.Second),
		storeFile:    persistencyFile,
		store:        lines,
		msgChan:      make(chan *parsing.Message),
		stop:         make(chan bool),
		storeMutex:   &sync.RWMutex{},
	}
	go memStore.listenForMsg()
	go memStore.runOldMessageCleaner()

	logging.Debug.Printf("Took me %s to load the file", time.Since(start))
	return memStore, nil
}

func (memStore *MemStore) runOldMessageCleaner() {
	logging.Debug.Println("Starting old Message Cleaner...")
	for {
		select {
		case <-memStore.stop:
			logging.Debug.Println("Stopping Cleanup Routine")
			return
		case <-memStore.cleanUpTick.C:
			logging.Debug.Println("run cleanup routine...")
			memStore.clean()
		}
	}
}

// Remove values to 80% of max sizes
func (memStore *MemStore) clean() {
	before_count := len(memStore.store)
	if before_count <= memStore.maximumLines {
		logging.Debug.Printf("have %d messages, no cleaning required", before_count)
		return
	}
	f, err := os.OpenFile(memStore.storeFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		logging.Error.Printf("Couldn't Open file for Reading before cleanup: %v", err)
		return
	}

	defer f.Close()

	memStore.storeMutex.Lock()

	target_count := int(float64(memStore.maximumLines) * AFTER_CLEAN_RATIO)
	logging.Debug.Printf("Have: %d Messages, Want %d. Cleaning...", before_count, target_count)

	removed_lines := memStore.store[0 : before_count-target_count]
	memStore.store = memStore.store[before_count-target_count : before_count]

	memStore.storeMutex.Unlock()

	// now append them out to file
	for _, line := range removed_lines {
		lineJson, jsonErr := json.Marshal(line)
		if jsonErr != nil {
			logging.Error.Printf("Error while Marshalling %v : %v", line, jsonErr)
		}
		_, writeErr := f.Write(lineJson)
		if writeErr != nil {
			logging.Error.Printf("Error while writing '%s' to '%s': %v", lineJson, memStore.storeFile, jsonErr)
		}
	}

}

func (memStore *MemStore) listenForMsg() {
	logging.Debug.Println("Starting Listen loop...")
	for {
		select {
		case <-memStore.stop:
			logging.Info.Println("Stopped listening for messages")
			return
		case msg := <-memStore.msgChan:
			memStore.appendMessage(msg)
		}
	}
}

func (memStore *MemStore) appendMessage(msg *parsing.Message) {
	memStore.storeMutex.Lock()
	defer memStore.storeMutex.Unlock()
	memStore.store = append(memStore.store, msg)
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
	memStore.storeMutex.RLock()
	defer memStore.storeMutex.RUnlock()

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
	memStore.cleanUpWaitGroup.Done()
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
	memStore.storeMutex.RLock()
	defer memStore.storeMutex.RUnlock()

	logging.Debug.Printf("Filtering for %s and %s with %d results ...", appFilter, msgFilter, num)
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
