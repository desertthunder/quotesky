package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/desertthunder/quotesky/lib/utils"
)

type Repository interface{}
type AppRepository struct {
	conn *sql.DB
	Log  *log.Logger
}

func InitAppRepo(dbg bool) *AppRepository {
	dbc := Connect(dbg)
	l := log.NewWithOptions(os.Stderr, utils.Options("App Repo 🗂️", dbg))
	return &AppRepository{dbc.db, l}
}

func (a AppRepository) CreateOrUpdate(h string, t string) error {
	id := 0

	err := a.conn.QueryRow(`SELECT id FROM apps WHERE handle = ?`, h).Scan(&id)

	if err == sql.ErrNoRows {
		a.Log.Debug("creating record")
	}

	if err != nil {
		return err
	}

	if id == 0 {
		now := time.Now().Format(time.RFC3339)
		err := a.conn.QueryRow(
			"INSERT INTO apps (handle, token, created_at, updated_at) "+
				"VALUES (?, ?, ?, ?) RETURNING id", h, t, now, now,
		).Scan(&id)

		if err != nil {
			return err
		}

		return nil
	}

	a.Log.Infof("app %d already exists, updating token", id)

	res, err := a.conn.Exec(
		`UPDATE apps SET handle = ?, token = ? WHERE id = ?`,
		h, t, id,
	)

	if err != nil {
		return err
	}

	af, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if af != 1 {
		return fmt.Errorf("%d was updated instead of %d", af, id)
	}

	return nil

}
