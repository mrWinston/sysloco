package sqlite

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mrWinston/sysloco/receiver/logging"
	"github.com/mrWinston/sysloco/receiver/parsing"
)

type SqliteStore struct {
	msgChan   chan *parsing.Message
	stop      chan bool
	storeFile string
	db        *sql.DB
}

const BUFFER_SIZE = 1000

func NewSqliteStore(storeFile string) (*SqliteStore, error) {
	logging.Info.Printf("Opening Sqlite DB at: %s\n", storeFile)

	db, err := sql.Open("sqlite3", storeFile)

	if err != nil {
		logging.Error.Printf("Error Opening Database file: %v", err)
		return nil, err
	}

	_, err = db.Exec(statementStrings.createLogTable)

	if err != nil {
		logging.Error.Printf("Error Creating Datbase: %v\n", err)
		return nil, err
	}

	var store = &SqliteStore{
		msgChan:   make(chan *parsing.Message, BUFFER_SIZE),
		stop:      make(chan bool),
		storeFile: storeFile,
		db:        db,
	}

	go store.listenForMsg()
	return store, nil
}

func (sqliteStore *SqliteStore) mustInsertMessage(msg *parsing.Message) {
	res, err := sqliteStore.db.Exec(
		statementStrings.insertLogEntry,
		msg.Priv,
		msg.Version,
		msg.Timestamp.Format(time.RFC3339Nano),
		msg.Hostname,
		msg.Appname,
		msg.Procid,
		msg.Msgid,
		msg.Msg,
	)
	if err != nil {
		logging.Error.Printf("Error Running Statement: %v, %v\n", res, err)
	}
}

func (sqliteStore *SqliteStore) listenForMsg() {
	logging.Debug.Println("Starting message listen loop")

	for {
		select {
		case <-sqliteStore.stop:
			logging.Debug.Println("Stop message listen loop")
			return
		case msg := <-sqliteStore.msgChan:
			sqliteStore.mustInsertMessage(msg)
		}
	}
}

// Store puts the parsing.Message into the memStore. Adding happens
// asyncronously, so calling GetLatest right after adding will likely not
// return the just added msg
func (sqliteStore *SqliteStore) Store(msg parsing.Message) error {
	sqliteStore.msgChan <- &msg
	logging.Debug.Printf("Added Msg, Current Store lenght: %v\n", len(sqliteStore.msgChan))
	return nil
}

// GetLatest Returns the Last added elements to the store, In Reverse order ( from
// newest to oldest )
func (sqliteStore *SqliteStore) GetLatest(number int) ([]*parsing.Message, error) {
	rows, err := sqliteStore.db.Query(
		statementStrings.getNLatest,
		number,
	)

	if err != nil {
		logging.Error.Printf("Error while fetching %d Rows from logs: %v\n", number, err)
		return nil, err
	}

	return get, nil
}
func (sqliteStore *SqliteStore) Filter(appRegex string, msgRegex string, num int) ([]*parsing.Message, error) {
	return nil, nil
}

func (sqliteStore *SqliteStore) Release() error {
	sqliteStore.stop <- true

	close(sqliteStore.msgChan)
	close(sqliteStore.stop)
	return sqliteStore.db.Close()

}
