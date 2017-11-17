package websocket

import (
	"io"
	"net/http"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/gobwas/ws"
)

var (
	DefaultHeaders = http.Header{}
)

// UpgradeHandler is a net/http.Handler which is responsible
// for incomming HTTP request upgrade process and handling the
// upgraded request with some Handler which works only with
// websockets.
type UpgradeHandler struct {
	Handler
	log loggers.Logger
}

func (h *UpgradeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()

	var (
		c   io.WriteCloser
		err error
	)

	c, _, _, err = ws.UpgradeHTTP(
		r,
		w,
		DefaultHeaders,
	)
	if err != nil {
		h.log.Error(err)
		return
	}

	h.ServeWebsocket(c, r)
}

func NewUpgradeHandler(h Handler, l loggers.Logger) *UpgradeHandler {
	return &UpgradeHandler{
		Handler: h,
		log: prefixwrapper.New(
			"UpgradeHandler: ",
			l,
		),
	}
}
