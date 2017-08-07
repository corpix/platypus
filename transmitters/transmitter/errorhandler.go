package transmitter

import (
	"io"
)

type ErrorHandler func(w io.Writer, err error)
