package api

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
const createPost string = "app.bsky.feed.post"

type SessionRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type Credentials struct {
	Handle       string
	Password     string
	AccessToken  string
	RefreshToken string
	DID          string
}

type Session struct {
	AccessJwt       string      `json:"accessJwt"`
	RefreshJwt      string      `json:"refreshJwt"`
	Handle          string      `json:"handle"`
	Did             string      `json:"did"`
	DidDoc          interface{} `json:"didDoc"`
	Email           string      `json:"email"`
	EmailConfirmed  bool        `json:"emailConfirmed"`
	EmailAuthFactor bool        `json:"emailAuthFactor"`
	Active          bool        `json:"active"`
	Status          string      `json:"status"`
}

type Client struct {
	Service     string
	Credentials *Credentials
}

func credentials() *Credentials {
	handle := os.Getenv("BLUESKY_HANDLE")
	password := os.Getenv("BLUESKY_PASSWORD")

	return &Credentials{
		Handle:   handle,
		Password: password,
	}
}

func (c *Credentials) SetSession(s Session) {
	c.AccessToken = s.AccessJwt
	c.RefreshToken = s.RefreshJwt
	c.DID = s.Did
}

func (c Client) BuildURL(service string, path string) string {
	return fmt.Sprintf("%s/xrpc/%s", service, path)
}

func (c Client) CreateSession() (*Session, error) {
	uri := c.BuildURL(c.Service, createSession)
	r := SessionRequest{c.Credentials.Handle, c.Credentials.Password}
	req, err := json.Marshal(r)

	if err != nil {
		log.Errorf("unable to build body: %s", err.Error())
		return nil, err
	}

	req_body := bytes.NewBuffer(req)
	rsp, err := http.Post(uri, "application/json", req_body)

	if err != nil {
		log.Errorf("unable to authenticate: %s", err.Error())
		return nil, err
	}

	defer rsp.Body.Close()

	rspBody, err := io.ReadAll(rsp.Body)

	if err != nil {
		log.Errorf("unable to read response %s", err.Error())
		return nil, err
	}

	s := Session{}

	json.Unmarshal(rspBody, &s)

	fmt_j, _ := json.MarshalIndent(s, "", "  ")

	log.Debugf(string(fmt_j))
	log.Infof("session created at %s", time.Now().Format("03:04 PM on 01/02/2006"))

	return &s, nil
}

func Init(s string) *Client {
	return &Client{Service: s, Credentials: credentials()}
}

func (c Client) CreatePost() {}
