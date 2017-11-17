package websocket

import (
	"io"
	"net/http"
)

type Handler interface {
	ServeWebsocket(w io.WriteCloser, r *http.Request)
}
