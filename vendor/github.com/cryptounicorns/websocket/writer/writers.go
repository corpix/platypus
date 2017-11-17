package writer

import (
	"io"

	"github.com/gobwas/ws/wsutil"
)

func wrapWsutil(f func(io.Writer, []byte) error) func(io.WriteCloser, []byte) error {
	return func(w io.WriteCloser, buf []byte) error {
		return f(w, buf)
	}
}

func NewServerBinary(w io.WriteCloser) io.WriteCloser {
	return New(w, wrapWsutil(wsutil.WriteServerBinary))
}

func NewServerText(w io.WriteCloser) io.WriteCloser {
	return New(w, wrapWsutil(wsutil.WriteServerText))
}

func NewClientBinary(w io.WriteCloser) io.WriteCloser {
	return New(w, wrapWsutil(wsutil.WriteClientBinary))
}

func NewClientText(w io.WriteCloser) io.WriteCloser {
	return New(w, wrapWsutil(wsutil.WriteClientText))
}
