package writers

import (
	"io"
)

type Writers interface {
	Iter(func(io.Writer))
}
