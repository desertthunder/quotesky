package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

const createSession string = "com.atproto.server.createSession"

type Credentials struct {
	Handle   string
	Password string
	JWT      string
}

type IClient interface {
	CreateSession()
	BuildURL(service string, path string) string
}

type Client struct {
	Service     string
	Credentials *Credentials
}

func SetCredentials() *Credentials {
	handle := os.Getenv("BLUESKY_HANDLE")
	password := os.Getenv("BLUESKY_PASSWORD")

	return &Credentials{
		Handle:   handle,
		Password: password,
	}
}

func (c *Credentials) SetJWT(tok string) {
	c.JWT = tok
}

func (c Client) BuildURL(service string, path string) string {
	return fmt.Sprintf("%s/xrpc/%s", service, path)
}

type createSessionRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func (c Client) CreateSession() error {
	uri := c.BuildURL(c.Service, createSession)
	r := createSessionRequest{
		Identifier: c.Credentials.Handle,
		Password:   c.Credentials.Password,
	}

	req, err := json.Marshal(r)

	if err != nil {
		log.Errorf("unable to build body: %s", err.Error())
		return err
	}

	req_body := bytes.NewBuffer(req)

	rsp, err := http.Post(uri, "application/json", req_body)

	if err != nil {
		log.Errorf("unable to authenticate: %s", err.Error())
		return err
	}

	defer rsp.Body.Close()

	rspBody, err := io.ReadAll(rsp.Body)

	if err != nil {
		log.Errorf("unable to read response %s", err.Error())
		return err
	}

	c.Credentials.SetJWT(string(rspBody))

	log.Infof("session created at %s", time.Now().Format("03:04 PM on 01/02/2006"))

	return nil
}
