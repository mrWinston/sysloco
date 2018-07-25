package storage

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/dgraph-io/badger"
	"github.com/mrWinston/sysloco/receiver/parsing"
)

// The BadgerStore provide the methods for interacting with the Badger Storage engine
type BadgerStore struct {
	db       *badger.DB
	options  badger.Options
	sequence *badger.Sequence
}

// BadgerDefaultOptions returns the Default Options struct from the badger package
func BadgerDefaultOptions() badger.Options {
	return badger.DefaultOptions
}

// NewBadgerStore creates a new BadgerStore from the Options provided. Returns
// an error if something goes wrong.
func NewBadgerStore(options badger.Options) (*BadgerStore, error) {
	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}
	seq, err := db.GetSequence([]byte("log"), 10)
	if err != nil {
		return nil, err
	}
	return &BadgerStore{
		db:       db,
		options:  options,
		sequence: seq,
	}, nil

}

// Store puts the parsing.Message struct in the Badger Storage.
func (badgerStore *BadgerStore) Store(message parsing.Message) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}
	key := strconv.FormatInt(message.Timestamp.Unix(), 10)

	badgerStore.db.Update(func(txn *badger.Txn) error {

		err := txn.Set([]byte(key), jsonMessage)
		return err

	})
	return nil
}

// GetLatest Returns the "number" last elements added to the Store
func (badgerStore *BadgerStore) GetLatest(number int) ([]*parsing.Message, error) {

	var results []*parsing.Message

	err := badgerStore.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = number
		opts.Reverse = false

		it := txn.NewIterator(opts)
		defer it.Close()

		var index = 0
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			val, err := item.Value()
			if err != nil {
				return err
			}
			var msg parsing.Message
			err = json.Unmarshal(val, &msg)
			if err != nil {
				return err
			}
			results[index] = &msg
			index++
			if index == number {
				break
			}
		}
		return nil
	})
	return results, err
}

// Release closes the Connection to the badger store. Should always be called
// before shutdown to ensure a clean exit
func (badgerStore *BadgerStore) Release() error {
	err := badgerStore.sequence.Release()
	if err != nil {
		return err
	}
	err = badgerStore.db.Close()
	if err != nil {
		return err
	}
	return nil
}

// Filter returns all elements in the store, whose Msg matches the provided
// regular expression. Elements are returned in reverse order ( from newest, to
// oldest )
func (badgerStore *BadgerStore) Filter(regex string) ([]*parsing.Message, error) {
	return nil, errors.New("Not Implemented Yet")
}
