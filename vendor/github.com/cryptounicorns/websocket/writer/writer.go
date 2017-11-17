package writer

import (
	"io"
)

type WriteWith = func(w io.WriteCloser, buf []byte) error

type Writer struct {
	io.WriteCloser
	writeWith WriteWith
}

func (w *Writer) Write(m []byte) (int, error) {
	var (
		err error
	)

	err = w.writeWith(
		w.WriteCloser,
		m,
	)
	if err != nil {
		return 0, err
	}

	return len(m), nil
}

func New(w io.WriteCloser, ww WriteWith) *Writer {
	return &Writer{
		WriteCloser: w,
		writeWith:   ww,
	}
}
