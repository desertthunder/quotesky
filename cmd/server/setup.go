package server

import (
	"github.com/charmbracelet/log"
	"github.com/desertthunder/quotesky/lib/api"
	"github.com/desertthunder/quotesky/lib/db"
	"github.com/desertthunder/quotesky/lib/utils"
	"github.com/urfave/cli/v2"
)

const dir string = "lib/db/migrations"
const service string = "https://bsky.social"
const env_path string = ".env"

func Setup() *cli.Command {
	return &cli.Command{
		Name:      "setup",
		Usage:     "create the local store & authenticate",
		UsageText: `Create the database.`,
		Aliases:   []string{"s"},
		Action: func(ctx *cli.Context) error {
			err := utils.LoadEnv(env_path)

			if err != nil {
				log.Errorf("unable to load env file: %s", err.Error())
				return err
			}

			dbc := db.Connect(true)
			r := db.Runner(dir, dbc, true)

			if err := r.Execute(); err != nil {
				dbc.Log.Errorf("unable to execute: %s", err.Error())
				return err
			}

			c := api.Init(service, true)
			a := db.InitAppRepo(true)
			s, err := c.CreateSession()

			if err != nil {
				return err
			}

			err = a.CreateOrUpdate(s.Handle, s.AccessJwt)

			if err != nil {
				return err
			}

			return nil
		},
	}
}
