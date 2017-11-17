package pooled

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/cryptounicorns/websocket/client/dialer"
	"github.com/cryptounicorns/websocket/client/request"
)

type Dialer struct {
	dialer.Dialer
}

func (d *Dialer) Dial(network string, addr string) (net.Conn, error) {
	// FIXME: You have inner dialer.
	return d.DialContext(
		context.Background(),
		network,
		addr,
	)
}

func (d *Dialer) DialContext(ctx context.Context, network string, addr string) (net.Conn, error) {
	return d.DialContext(ctx, network, addr)
}

func New() *Dialer {
	// FIXME: create inner dialer?
	return &Dialer{}
}
