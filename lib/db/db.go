// Database connection
package db

import (
	"database/sql"
	"os"
	"time"

	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"
)

type DBConn struct {
	db  *sql.DB
	Log *log.Logger
}

// Create/Connect to database
func Connect(dbg bool) *DBConn {
	var err error

	opts := log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "Database 💾",
	}

	if dbg {
		opts.Level = log.DebugLevel
	}

	conn := DBConn{Log: log.NewWithOptions(os.Stderr, opts)}
	conn.db, err = sql.Open("sqlite3", "db.sqlite3")

	if err != nil {
		conn.Log.Errorf(
			"unable to connect to database: %s", err.Error(),
		)

		os.Exit(1)
	}

	return &conn
}
