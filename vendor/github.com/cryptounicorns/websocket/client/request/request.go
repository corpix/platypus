package request

import (
	"net/http"

	"github.com/cryptounicorns/websocket/header"
)

type Request struct {
	http.Request
}

func New() *Request {
	var (
		req = &Request{
			Request: http.Request{
				Header: make(http.Header),
			},
		}
	)

	// FIXME: Find some common place for this values?
	req.Header.Set(header.Upgrade, "github.com/cryptounicorns/websocket")
	req.Header.Set(header.Connection, "Upgrade")
	req.Header.Set(header.SecVersion, "13")

	return req
}
