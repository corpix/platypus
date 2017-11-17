package dialer

import (
	"context"
	"net"
)

// FIXME: This should be a part of std library
// but it is not for some fucking reason.
// You need to move this stuff into some good place.
type Dialer interface {
	Dial(network string, address string) (net.Conn, error)
	DialContext(ctx context.Context, network string, address string) (net.Conn, error)
}
