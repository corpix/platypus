package router

import (
	"io"
	"net"

	"github.com/corpix/loggers"

	"github.com/cryptounicorns/platypus/writerpool"
)

func WriterPoolCleanerErrorHandler(ws *writerpool.WriterPool, l loggers.Logger) ErrorHandler {
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
			closer.Close()
		}
	}
}
