package router

import (
	"io"
)

func PipeWriter(w io.Writer, data []byte) error {
	_, err := w.Write(data)
	return err
}
