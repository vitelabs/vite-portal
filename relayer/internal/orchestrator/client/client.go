package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vitelabs/vite-portal/shared/pkg/crypto"
	"github.com/vitelabs/vite-portal/shared/pkg/rpc"
	sharedtypes "github.com/vitelabs/vite-portal/shared/pkg/types"
)

type Client struct {
	url        string
	timeout    time.Duration
	jwtHandler crypto.JWTHandler
	ws         *rpc.Client
}

func NewClient(url, jwtSecret string, timeout, jwtExpiryTimeout time.Duration) *Client {
	return &Client{
		url:        url,
		timeout:    timeout,
		jwtHandler: *crypto.NewJWTHandler([]byte(jwtSecret), jwtExpiryTimeout),
	}
}

func (c *Client) Connect(jwtSubject string) error {
	if c.ws != nil {
		c.ws.Close()
	}
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: c.timeout,
	}
	token := c.CreateToken(jwtSubject, 10) // expires in 10 seconds
	client, err := rpc.DialWebsocketWithDialer(context.Background(), c.url, "", fmt.Sprintf("Bearer %s", token), dialer)
	if err != nil {
		return err
	}
	c.ws = client
	return nil
}

func (c *Client) CreateToken(jwtSubject string, exp int64) string {
	return c.jwtHandler.IssueDefaultToken(jwtSubject, sharedtypes.JWTRelayerIssuer, exp)
}

func (c *Client) Close() {
	c.ws.Close()
}

func (c *Client) RegisterNames(apis []rpc.API) error {
	for _, api := range apis {
		if err := c.ws.RegisterName(api.Namespace, api.Service); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Closed() <-chan interface{} {
	return c.ws.WriteConn.Closed()
}

func (c *Client) Call(result interface{}, method string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.ws.CallContext(ctx, result, method, args...)
}
