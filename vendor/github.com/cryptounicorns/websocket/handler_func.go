package websocket

import (
	"io"
	"net/http"
)

type HandlerFunc func(w io.Writer, r *http.Request)
