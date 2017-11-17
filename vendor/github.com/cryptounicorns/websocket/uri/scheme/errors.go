package scheme

import (
	"fmt"
)

type ErrUnsupportedScheme struct {
	Scheme string
}

func (e *ErrUnsupportedScheme) Error() string {
	return fmt.Sprintf(
		"Unsupported scheme: %s",
		e.Scheme,
	)
}

func NewErrUnsupportedScheme(scheme string) *ErrUnsupportedScheme {
	return &ErrUnsupportedScheme{
		Scheme: scheme,
	}
}
