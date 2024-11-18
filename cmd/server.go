package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/desertthunder/quotesky/lib"
	"github.com/urfave/cli/v2"
)

func handleConnError(e error, c net.Conn) {
	out := fmt.Sprintf("something went wrong: %s", e.Error())

	log.Error(out)
	c.Write([]byte(out))
}
func Handler(conn net.Conn) {
	defer conn.Close()

	log.Infof("Handling connection from %s\n", conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadString('\n')

		if err != nil && err == io.EOF {
			continue
		}

		if err != nil {
			handleConnError(err, conn)
			continue
		}

		msg := Message{}
		err = json.Unmarshal([]byte(data), &msg)

		if err != nil {
			handleConnError(err, conn)
			continue
		}

		s := fmt.Sprintf("content: %s", msg.Content)
		if len(msg.Hashtags) > 0 {
			s = fmt.Sprintf("%s | hashtags: %s", s, strings.Join(msg.Hashtags, ", "))
		}

		log.Info(s)
		_, err = conn.Write([]byte("OK\n\n"))

		if err != nil {
			handleConnError(err, conn)
		}
	}
}

func Beat(d time.Duration) {
	t := time.NewTicker(d)

	defer t.Stop()

	for {
		select {
		case tick := <-t.C:
			log.Infof("heartbeat at %s on %s",
				tick.Format("03:04:05 PM"), tick.Format("01/02/2006"))
		}
	}
}

func Server(p string, b int) error {
	addr := fmt.Sprintf(":%s", p)
	l, err := net.Listen("tcp", addr)

	if err != nil {
		log.Errorf("unable to open listener: %s", err.Error())

		return err
	}

	defer l.Close()

	log.Infof("listening at %s", p)

	go Beat(time.Duration(b) * time.Second)

	for {

		conn, err := l.Accept()

		if err != nil {
			log.Errorf("unable to accept messages: %s", err.Error())
			return err
		}

		go Handler(conn)
	}
}

func RunServer(p int) *cli.Command {
	lib.LoadEnv(env_path)

	// client := lib.Client{
	// 	Service:     service,
	// 	Credentials: lib.SetCredentials(),
	// }

	return &cli.Command{
		Name:    "tcp",
		Usage:   "run the tcp server",
		Aliases: []string{"t"},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "port",
				Value: p,
			},
			&cli.IntFlag{
				Name:  "beat",
				Value: 2,
			},
		},
		Action: func(ctx *cli.Context) error {
			port := strconv.Itoa(ctx.Int("port"))
			beat := ctx.Int("beat")
			Server(port, beat)
			return nil
		},
	}
}
