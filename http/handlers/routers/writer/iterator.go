package writer

import (
	"io"
)

type Iterator interface {
	Iter(func(io.Writer))
}
