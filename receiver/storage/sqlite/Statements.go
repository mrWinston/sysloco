package sqlite

var statementStrings = struct {
	createLogTable     string
	insertLogEntry     string
	getNLatest         string
	getNLatestFiltered string
}{
	createLogTable: `
		CREATE TABLE IF NOT EXISTS logentries
		(
			id INTEGER PRIMARY KEY ASC,
			priv VARCHAR(127),
			version INTEGER,
			timestamp DATETIME,
			hostname VARCHAR(255),
			appname VARCHAR(255),
			procid VARCHAR(255),
			msgid VARCHAR(255),
			msg MEDIUMTEXT
		);`,
	insertLogEntry: `
		INSERT INTO logentries (
			priv, version, timestamp, hostname, appname, procid, msgid, msg
		) VALUES(
			?, ?, ?, ?, ?, ?, ?, ?
		);
	`,
	getNLatest: `
		SELECT
			*
		FROM
			logentries
		ORDER BY
			id DESC
		LIMIT ?;
	`,
	getNLatestFiltered: `
		SELECT
			*
		FROM
			logentries
		ORDER BY
			id DESC
		WHERE
			appname regexp ? AND
			msg regexp ?
		LIMIT ?;
	`,
}
