package routers

import (
	"io"
	"net"

	"github.com/corpix/loggers"

	"github.com/cryptounicorns/platypus/http/handlers/routers/errors"
	"github.com/cryptounicorns/platypus/iopool"
)

func NewWriterPoolCleaner(ws *iopool.Writer, l loggers.Logger) errors.Handler {
	return func(w io.Writer, err error) {
		var (
			closer io.Closer
			ok     bool
		)

		ws.Remove(w)

		_, ok = err.(*net.OpError)
		if !ok {
			l.Error(err)
		}

		closer, ok = w.(io.Closer)
		if ok {
			_ = closer.Close()
		}
	}
}
