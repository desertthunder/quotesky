package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/desertthunder/quotesky/lib/api"
	"github.com/desertthunder/quotesky/lib/utils"
	"github.com/urfave/cli/v2"
)

type Protocol struct {
	client   *api.Client
	port     int
	addr     string
	beat     time.Duration
	logger   *log.Logger
	conn     net.Conn
	listener net.Listener
}

// Writes response to connection
func (p Protocol) handleConnError(e error) {
	out := fmt.Sprintf("something went wrong: %s", e.Error())

	p.logger.Error(out)

	_, err := p.conn.Write([]byte(out))

	if err != nil {
		p.logger.Errorf("unable to write to connection %s", err.Error())
	}
}

func (p Protocol) handleMessage() {
	defer p.conn.Close()

	log.Infof("Handling connection from %s\n", p.conn.RemoteAddr().String())
	reader := bufio.NewReader(p.conn)

	for {
		data, err := reader.ReadString('\n')

		if err != nil && err == io.EOF {
			continue
		}

		if err != nil {
			p.handleConnError(err)
			continue
		}

		msg := api.Message{}
		err = json.Unmarshal([]byte(data), &msg)

		if err != nil {
			p.handleConnError(err)
			continue
		}

		s := fmt.Sprintf("content: %s", msg.Content)
		if len(msg.Hashtags) > 0 {
			s = fmt.Sprintf("%s | hashtags: %s", s, strings.Join(msg.Hashtags, ", "))
		}

		log.Info(s)

		err = p.client.CreatePost(msg)

		if err != nil {
			p.handleConnError(err)
		}

		_, err = p.conn.Write([]byte("OK\n\n"))

		if err != nil {
			p.handleConnError(err)
		}
	}
}

func (p Protocol) heartbeat() {
	t := time.NewTicker(p.beat)

	defer t.Stop()

	for {
		tick := <-t.C

		log.Infof("heartbeat at %s on %s", tick.Format("03:04:05 PM"), tick.Format("01/02/2006"))
	}
}

func (p *Protocol) updateConn(c net.Conn) {
	p.conn = c
}

func (p Protocol) listen() error {

	defer p.listener.Close()

	log.Infof("listening at %s", p.addr)

	go p.heartbeat()

	for {
		conn, err := p.listener.Accept()

		p.updateConn(conn)

		if err != nil {
			log.Errorf("unable to accept messages: %s", err.Error())
			return err
		}

		go p.handleMessage()
	}
}

func (p *Protocol) SetListener() {
	var err error
	p.listener, err = net.Listen("tcp", p.addr)

	if err != nil {
		p.logger.Errorf("unable to open listener: %s", err.Error())
		os.Exit(1)
	}
}

// Set protocol heart rate in seconds
func (p *Protocol) SetHeartRate(hr int) {
	p.beat = time.Duration(hr) * time.Second
}

// Sets listener address to port pt
func (p *Protocol) SetAddress(pt int) {
	p.port = pt
	p.addr = fmt.Sprintf(":%s", strconv.Itoa(p.port))
}

// Sets prefixed logger to protocol process
func (p *Protocol) SetLogger(opts *log.Options) {
	if opts == nil {
		p.logger = log.NewWithOptions(
			os.Stderr,
			log.Options{
				Level:           log.DebugLevel,
				ReportCaller:    true,
				ReportTimestamp: true,
				TimeFormat:      time.Stamp,
				Prefix:          "Protocol 🖧",
			},
		)

		return
	}

	p.logger = log.NewWithOptions(os.Stderr, *opts)
}

func (p *Protocol) SetClient() {
	p.client = api.Init(service, true)
	_, err := p.client.CreateSession()

	if err != nil {
		log.Errorf("unable to set session: %s", err.Error())
	}
}

// Protocol constructor
func protocol(p int, b int) *Protocol {
	pr := Protocol{}
	pr.SetAddress(p)
	pr.SetHeartRate(b)
	pr.SetListener()
	pr.SetLogger(nil)
	pr.SetClient()

	return &pr
}

// Run tcp listener
func run(ctx *cli.Context) error {
	if ctx.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	err := utils.LoadEnv(env_path)

	if err != nil {
		return err
	}

	port := ctx.Int("port")
	beat := ctx.Int("beat")

	p := protocol(port, beat)

	if err := p.listen(); err != nil {
		log.Errorf("protocol issue: %s", err.Error())
		return err
	}

	return nil
}

// server Command definition
func RunServer(p int) *cli.Command {
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
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
			},
		},
		Action: run,
	}
}
