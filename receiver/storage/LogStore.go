package storage

import "github.com/mrWinston/sysloco/receiver/parsing"

// LogStore provides a Common interface for all Structs that act as a Logging
// Storage Solution.
type LogStore interface {

	// Store the parsing.Message in a LogStore.  Returns an
	// error in case sth goes south
	Store(msg parsing.Message) error

	// Return the last *number* stored items in the store
	GetLatest(number int) ([]*parsing.Message, error)

	// Release all resources held by this Storage engine and close all pending
	// connections
	Release() error

	// Run this Filter through the DB. Error should not be set, if the filter
	// didn't return any results. Instead, an Empty slice shall be returned.
	// it filters the app and Msg field of the logmsg and returns at most num
	// results
	Filter(appRegex string, msgRegex string, num int) ([]*parsing.Message, error)
}
