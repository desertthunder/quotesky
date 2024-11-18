package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"
)

const env_path string = ".env"
const service string = "https://bsky.social"
const Port int = 9000

var commands []*cli.Command

type Message struct {
	Content  string
	Hashtags []string
}

func Execute(p int) error {
	app := &cli.App{
		Name:  "qsky",
		Usage: "make posts to bluesky",
		Action: func(*cli.Context) error {
			log.Info("execute quotesky")
			return nil
		},
		Commands: []*cli.Command{RunServer(p), Post()},
	}

	return app.Run(os.Args)
}
