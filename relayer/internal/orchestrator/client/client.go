package client

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	url string
	timeout time.Duration
	Conn *websocket.Conn
}

func NewClient(url string, timeout int64) *Client {
	return &Client{
		url: url,
		timeout: time.Duration(timeout) * time.Millisecond,
	}
}

func (c *Client) Connect() error {
	dialer := websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		HandshakeTimeout: c.timeout,
	}
	conn, _, err := dialer.Dial(c.url, nil)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}