package sqlite

import (
	"database/sql"
	"time"

	"github.com/mrWinston/sysloco/receiver/parsing"
)

func getMessageFromRows(rows *sql.Rows) ([]*parsing.Message, error) {
	var res []*parsing.Message = make([]*parsing.Message, 0)

	for rows.Next() {
		var (
			id        int
			priv      string
			version   int
			timestamp time.Time
			hostname  string
			appname   string
			procid    string
			msgid     string
			msg       string
		)

		err := rows.Scan(
			&id,
			&priv,
			&version,
			&timestamp,
			&hostname,
			&appname,
			&procid,
			&msgid,
			&msg,
		)

		if err != nil {
			return nil, err
		}

		res = append(res, &parsing.Message{
			Priv:      priv,
			Version:   version,
			Timestamp: timestamp,
			Hostname:  hostname,
			Appname:   appname,
			Procid:    procid,
			Msgid:     msgid,
			Msg:       msg,
		})

	}

	return res, nil

}
