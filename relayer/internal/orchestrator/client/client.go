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

func NewClient(url string, timeout time.Duration) *Client {
	return &Client{
		url: url,
		timeout: timeout,
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