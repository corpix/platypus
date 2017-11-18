package errors

import (
	"io"
)

type Handler = func(w io.Writer, err error)
