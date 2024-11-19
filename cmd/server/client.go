package server

import (
	"github.com/desertthunder/quotesky/lib/db"
	"github.com/urfave/cli/v2"
)

const dir string = "lib/db/migrations"

func Setup() *cli.Command {
	return &cli.Command{
		Name:      "setup",
		Usage:     "create the local store & authenticate",
		UsageText: `Create the database.`,
		Aliases:   []string{"s"},
		Action: func(ctx *cli.Context) error {
			dbc := db.Connect(true)
			r := db.Runner(dir, dbc, true)

			if err := r.Execute(); err != nil {
				dbc.Log.Errorf("unable to execute: %s", err.Error())
				return err
			}

			return nil
		},
	}
}
