// Database connection
package db

import (
	"database/sql"
	"os"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

type DBConn struct {
	db     *sql.DB
	level  log.Level
	logger *log.Logger
}

func Connect(dbg bool) *DBConn {
	var err error
	conn := DBConn{level: log.InfoLevel}

	if dbg {
		conn.level = log.DebugLevel
	}

	conn.db, err = sql.Open("sqlite", "db.sqlite3")

	if err != nil {
		conn.logger.Errorf("unable to connect to database: %s", err.Error())

		os.Exit(1)
	}

	return &conn
}
