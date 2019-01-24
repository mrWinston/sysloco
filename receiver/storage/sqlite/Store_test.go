package sqlite

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestNewSqliteStore(t *testing.T) {

	_, err := NewSqliteStore("/tmp/testdb.sqlite3")

	if err != nil {
		t.Fail()
	}

}
