package reader

import (
	"io"
	"io/ioutil"

	"github.com/gobwas/ws"

	"github.com/cryptounicorns/websocket/capability"
	"github.com/cryptounicorns/websocket/reader/frame"
)

type Reader struct {
	Source     io.Reader
	Capability capability.Capability

	OnContinuation FrameHandler
	OnIntermediate FrameHandler

	frame *frame.Frame
}

func New(r io.Reader, c capability.Capability) *Reader {
	return &Reader{
		Source:     r,
		Capability: c,
	}
}

func (r *Reader) Read(p []byte) (int, error) {
	var (
		n   int
		err error
	)

	if r.frame == nil {
		r.frame, err = frame.Next(
			r.Source,
			r.Capability,
		)
		if err != nil {
			return n, err
		}
	}

	n, err = r.frame.Read(p)
	if err == io.EOF {
		r.frame = nil
		return n, nil
	}

	return n, err
}
