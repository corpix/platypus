package websocket

import (
	"io"
	"net/http"
)

type FuncHandler struct {
	handlerFunc HandlerFunc
}

func (h *FuncHandler) ServeWebsocket(w io.Writer, r *http.Request) {
	h.handlerFunc(w, r)
}

func NewFuncHandler(h HandlerFunc) *FuncHandler {
	return &FuncHandler{
		handlerFunc: h,
	}
}

func HandlerFuncWithWriter(h HandlerFunc, wc func(io.Writer) io.Writer) HandlerFunc {
	return func(w io.Writer, r *http.Request) {
		h(wc(w), r)
	}
}
