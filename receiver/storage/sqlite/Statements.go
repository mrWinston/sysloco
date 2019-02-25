package sqlite

var statementStrings = struct {
	createLogTable      string
	getNLatest          string
	getNLatestFilterApp string
	getNLatestFilterMsg string
	getNLatestFiltered  string
	insertLogEntry      string
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
		WHERE
			appname REGEXP ? AND
			msg REGEXP ?
		ORDER BY
			id DESC
		LIMIT ?;
	`,
	getNLatestFilterApp: `
		SELECT
			*
		FROM
			logentries
		WHERE
			appname REGEXP ?
		ORDER BY
			id DESC
		LIMIT ?;
	`,
	getNLatestFilterMsg: `
		SELECT
			*
		FROM
			logentries
		WHERE
			msg REGEXP ?
		ORDER BY
			id DESC
		LIMIT ?;
	`,
}
