// Migration execution
package db

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

// Migration slices
type Files struct {
	up   []fs.FileInfo
	down []fs.FileInfo
}

// Execution Context for Migrations
//
// Holds state and defines container for migration files
type MigrationRunner struct {
	Dir   string
	Conn  *DBConn
	Files Files
	Log   *log.Logger
}

// Migration Record
type Migration struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Applied   bool      `db:"applied"`
	CreatedAt time.Time `db:"created_at"`
}

// Build file path for migration file
//
// Uses dir from MigrationRunner
func (m Migration) GetPath(d string) string {
	return fmt.Sprintf("%s/%s", d, m.Name)
}

// Opens migrations directory and returns list of files
//
// Defaults to lib/db/migrations/*.sql
func (r *MigrationRunner) GetFiles() error {
	entries, err := os.ReadDir(r.Dir)

	if err != nil {
		return err
	}

	r.Files = Files{}

	for _, entry := range entries {
		info, err := entry.Info()

		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), "up.sql") {
			r.Files.up = append(r.Files.up, info)
		}

		if strings.HasSuffix(info.Name(), "down.sql") {
			r.Files.down = append(r.Files.down, info)
		}
	}

	return nil
}

// Creates Migration Records
//
// Inserts Rows into Database
func (r MigrationRunner) InsertMigrations() error {
	for _, f := range r.Files.up {
		m := Migration{Name: f.Name(), Applied: false}

		res, err := r.Conn.db.Exec(
			`INSERT INTO schema_migrations (name, applied) VALUES (?, ?) RETURNING id`,
			m.Name, m.Applied,
		)

		if err != nil {
			return err
		}

		id, err := res.LastInsertId()

		if err != nil {
			return err
		}

		r.Log.Debugf("Inserted migration: %d", id)
	}

	return nil
}

// Check for the existence of Migration Table
//
// Create it if it doesn't.
func (r MigrationRunner) CheckMigrationsTable() error {
	m := Migration{Name: "0000_init.sql"}
	f, err := os.ReadFile(m.GetPath(r.Dir))

	if err != nil {
		return err
	}

	query := string(f)
	_, err = r.Conn.db.Exec(query)

	return err
}

// Check migration record
//
// If applied, ensure that the forward migrations aren't reverted
func (r MigrationRunner) CheckApplied(f string) (bool, bool) {
	count := 0
	err := r.Conn.db.QueryRow(
		`SELECT COUNT(*) FROM schema_migrations WHERE name = ?`,
		f,
	).Scan(&count)

	if err != nil {
		r.Log.Errorf("unable to count migrations: %s", err.Error())
		return false, false
	}

	created := count > 0

	if created {
		m := Migration{}
		err := r.Conn.db.QueryRow(
			`SELECT * FROM schema_migrations WHERE name = ?`,
			f,
		).Scan(&m)

		if err != nil {
			r.Log.Errorf("unable to find migrations: %s", err.Error())
			return created, false
		}

		return created, m.Applied
	}

	return created, false
}

// Revert a specific migration
//
// Checks for the presence of a down migration file
func (r MigrationRunner) RevertMigration(mn string) {}

// Apply a specific migration
//
// Finds the row, runs the script, updates the row
func (r MigrationRunner) ApplyMigration(mn string) {}

// Runs migrations
func (r MigrationRunner) RunMigrations(d string) error {
	err := r.GetFiles()

	if err != nil {
		r.Log.Errorf("unable to get file list: %s", err.Error())
		os.Exit(1)
	}

	for _, f := range r.Files.up {
		name := strings.TrimSuffix(f.Name(), "up.sql")

		res, err := r.Conn.db.Exec(
			`UPDATE schema_migrations SET applied = TRUE WHERE name = ? RETURNING *`,
			f.Name(),
		)

		if err != nil {
			return err
		}

		n, err := res.RowsAffected()

		if err != nil {
			return err
		}

		if n != 1 {
			return fmt.Errorf("tried to update migration %s but updated %d rows", name, n)
		}
	}

	return nil
}

// MigrationRunner constructor
func Runner(d string, c *DBConn, dbg bool) *MigrationRunner {
	opts := log.Options{
		ReportTimestamp: true,
		ReportCaller:    true,
		TimeFormat:      time.Kitchen,
		Prefix:          "Runner 🚀",
	}

	if dbg {
		opts.Level = log.DebugLevel
	}

	return &MigrationRunner{Dir: d, Conn: c, Files: Files{}, Log: log.NewWithOptions(
		os.Stderr, opts,
	)}
}

// Checks for stored but unapplied migrations
func CheckPending() {}

// Runs migration process
func (r MigrationRunner) Execute() error {
	r.Log.Debug("starting execution")

	defer r.Conn.db.Close()

	err := r.CheckMigrationsTable()

	if err != nil {
		r.Log.Errorf("something went wrong: %s", err.Error())
		return err
	}

	err = r.GetFiles()

	if err != nil {
		r.Log.Errorf("something went wrong: %s", err.Error())
		return err
	}

	err = r.InsertMigrations()

	if err != nil {
		r.Log.Errorf("something went wrong: %s", err.Error())
		return err
	}

	return nil
}
