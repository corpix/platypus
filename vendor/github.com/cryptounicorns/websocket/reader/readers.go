package reader

import (
	"io"

	"github.com/gobwas/ws/wsutil"
)

func NewServerBinary(w io.Reader) io.Reader {
	return NewReader(w, wsutil.ReadServerBinary)
}

func NewServerText(w io.Reader) io.Reader {
	return NewReader(w, wsutil.ReadServerText)
}

func NewClientBinary(w io.Reader) io.Reader {
	return NewReader(w, wsutil.ReadClientBinary)
}

func NewClientText(w io.Reader) io.Reader {
	return NewReader(w, wsutil.ReadClientText)
}
