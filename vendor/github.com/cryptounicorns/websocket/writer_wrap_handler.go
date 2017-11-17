package websocket

import (
	"io"
	"net/http"
)

// WriterWrapHandler is a Handler which wraps
// some Handler with a io.Writer constructor.
type WriterWrapHandler struct {
	handler Handler

	// writer will receive a original io.Writer `a` and
	// must return a new io.Writer `b` which is usualy just
	// a wrapped `a`.
	writer func(a io.WriteCloser) (b io.WriteCloser)
}

func (h *WriterWrapHandler) ServeWebsocket(w io.WriteCloser, r *http.Request) {
	h.handler.ServeWebsocket(
		h.writer(w),
		r,
	)
}

func NewWriterWrapHandler(h Handler, wc func(io.WriteCloser) io.WriteCloser) *WriterWrapHandler {
	return &WriterWrapHandler{
		handler: h,
		writer:  wc,
	}
}
