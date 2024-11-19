package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/desertthunder/quotesky/lib/api"
	"github.com/urfave/cli/v2"
)

func Post() *cli.Command {
	return &cli.Command{
		Name:    "post",
		Aliases: []string{"p"},
		Usage:   "post to bluesky",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "content",
				Aliases:  []string{"c"},
				Usage:    "the body of your post",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     "hashtags",
				Aliases:  []string{"tags", "t"},
				Usage:    "any hashtags you want to add",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			log.Info("Making request to tcp client")

			content := ctx.String("content")
			hashtags := ctx.StringSlice("hashtags")

			if len(hashtags) > 0 {
				for i, h := range hashtags {
					hashtags[i] = fmt.Sprintf("#%s", h)
				}
			}

			log.Infof(
				"Sent message content:%s\nhashtags:%s",
				content, strings.Join(hashtags, ", "),
			)

			conn, err := net.Dial("tcp", ":9000")

			if err != nil {
				return err
			}

			msg := api.Message{Content: content, Hashtags: hashtags}
			data, err := json.Marshal(msg)

			if err != nil {
				return err
			}

			_, err = conn.Write(append(data, '\n'))
			reader := bufio.NewReader(conn)
			resp, err := reader.ReadString('\n')

			if err != nil {
				return err
			}

			log.Infof("response: %s", resp)

			return nil
		},
	}
}
