package writers

import (
	"io"
)

type Writer func(w io.Writer, data []byte) error
