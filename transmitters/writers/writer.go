package writers

import (
	"io"
)

type Writer func(io.Writer, []byte) error
