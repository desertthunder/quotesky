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

const dir string = "lib/db/migrations"

type files struct {
	up   []fs.FileInfo
	down []fs.FileInfo
}

// Execution Context for Migrations
//
// Holds state and defines container for migration files
type MigrationRunner struct {
	dir   string
	conn  DBConn
	files files
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
func (r MigrationRunner) GetFiles(d string) error {
	entries, err := os.ReadDir(d)

	if err != nil {
		return err
	}

	r.files = files{}

	for _, entry := range entries {
		info, err := entry.Info()

		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), "up.sql") {
			r.files.up = append(r.files.up, info)
		}

		if strings.HasSuffix(info.Name(), "down.sql") {
			r.files.down = append(r.files.down, info)
		}
	}

	return nil
}

// Creates Migration Records
//
// Inserts Rows into Database
func (r MigrationRunner) InsertMigrations() error {
	for _, f := range r.files.up {
		m := Migration{Name: f.Name(), Applied: false}
		res, err := r.conn.db.Exec(
			`INSERT INTO schema_migrations (filename, applied) VALUES (?, ?) RETURNING id`,
			m.Name, m.Applied,
		)

		if err != nil {
			return err
		}

		id, err := res.LastInsertId()

		if err != nil {
			return err
		}

		log.Debugf("Inserted migration: %d", id)
	}

	return nil
}

// Check for the existence of Migration Table
//
// Create it if it doesn't.
func (r MigrationRunner) CheckMigrationsTable() error {
	m := Migration{Name: "0000_init.sql"}
	f, err := os.ReadFile(m.GetPath(r.dir))

	if err != nil {
		return err
	}

	query := string(f)
	_, err = r.conn.db.Exec(query)

	return err
}

// Check migration record
//
// If applied, ensure that the forward migrations aren't reverted
func (r MigrationRunner) CheckApplied(f string) (bool, bool) {
	count := 0
	err := r.conn.db.QueryRow(
		`SELECT COUNT(*) FROM schema_migrations WHERE filename = ?`,
		f,
	).Scan(&count)

	if err != nil {
		log.Errorf("unable to count migrations: %s", err.Error())
		return false, false
	}

	created := count > 0

	if created {
		m := Migration{}
		err := r.conn.db.QueryRow(
			`SELECT * FROM schema_migrations WHERE filename = ?`,
			f,
		).Scan(&m)

		if err != nil {
			log.Errorf("unable to find migrations: %s", err.Error())
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
	err := r.GetFiles(d)

	if err != nil {
		log.Error("unable to get file list: %s", err.Error())
		os.Exit(1)
	}

	for _, f := range r.files.up {
		name := strings.TrimSuffix(f.Name(), "up.sql")

		res, err := r.conn.db.Exec(
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

// Checks for stored but unapplied migrations
func CheckPending() {}

// Runs migration process
func Execute()
