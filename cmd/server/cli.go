package server

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

const env_path string = ".env"
const service string = "https://bsky.social"
const Port int = 9000

func Execute(p int) error {
	app := &cli.App{
		Name:  "qsky",
		Usage: "make posts to bluesky",
		Action: func(*cli.Context) error {
			log.Info("execute quotesky")
			return nil
		},
		Commands: []*cli.Command{RunServer(p), Post(), Setup()},
	}

	return app.Run(os.Args)
}
