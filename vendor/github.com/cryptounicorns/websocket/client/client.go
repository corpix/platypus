package client

import (
	"context"
	"net"

	"github.com/cryptounicorns/websocket/client/dialer"
)

type Client struct {
	dialer.Dialer
}

func (c *Client) Handshake() {}

// FIXME: This naming is not so good
func (c *Client) Connect(network string, address string) (net.Conn, error) {
	return c.Dialer.Dial(network, address)
}

func (c *Client) ConnectContext(ctx context.Context, network string, address string) (net.Conn, error) {
	return c.Dialer.DialContext(ctx, network, address)
}

func New(d dialer.Dialer) *Client {
	return &Client{
		Dialer: d,
	}
}
