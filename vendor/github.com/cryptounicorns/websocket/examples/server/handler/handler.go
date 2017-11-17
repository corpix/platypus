package handler

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/cryptounicorns/websocket/examples/server/message"
)

type Handler struct {
	log loggers.Logger
}

func (h *Handler) ServeWebsocket(w io.WriteCloser, r *http.Request) {
	var (
		m   []byte
		err error
	)

	for {
		m, err = json.Marshal(
			message.New("hello"),
		)
		if err != nil {
			h.log.Error(err)
			return
		}

		_, err = w.Write(m)
		if err != nil {
			switch err.(type) {
			case *net.OpError:
			default:
				h.log.Error(err)
			}

			h.log.Debugf(
				"Terminating loop for %#v",
				w,
			)
			return
		}

		time.Sleep(1 * time.Second)
	}
}

func New(l loggers.Logger) *Handler {
	return &Handler{
		log: prefixwrapper.New("Handler: ", l),
	}
}
