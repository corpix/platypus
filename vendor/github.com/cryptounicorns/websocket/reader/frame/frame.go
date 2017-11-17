package frame

import (
	"io"

	"github.com/cryptounicorns/websocket/header"
	"github.com/gobwas/ws"
)

type Frame struct {
	r io.Reader
	io.Reader
}

func (f *Frame) Read(buf []byte) (int, error) {
	var (
		fh  ws.Header
		err error
	)

	// FIXME: Continue refactoring
	if f.r == nil {
		fh, err = ws.ReadHeader(f.Reader)
		if err != nil {
			return 0, err
		}

		err = header.Check(fh, c)
		if err != nil {
			// FIXME: Probably it is correct to return
			// The number of bytes we had read from f.Reader,
			// but we just can't obtain them yet.
			return 0, err
		}

		if fh.Masked {
			f.r = NewCipherReader(
				io.LimitReader(
					f.Reader,
					fh.Length,
				),
				fh.Mask,
			)
		} else {
			f.r = io.LimitReader(
				f.Reader,
				fh.Length,
			)
		}
	}

	return f.r.Read(buf)
}

func New(r io.Reader) *Frame {
	return &Frame{
		Reader: r,
	}
}
