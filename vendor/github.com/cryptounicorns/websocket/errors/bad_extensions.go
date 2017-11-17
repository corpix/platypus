package errors

import (
	"fmt"
	"strings"
)

type ErrBadExtensions struct {
	Extensions []string
}

func (e *ErrBadExtensions) Error() string {
	return fmt.Printf(
		"Bad sub-protocol: %s",
		strings.Join(e.Extensions, ", "),
	)
}

func NewErrBadExtensions(extensions []string) *ErrBadExtensions {
	return &ErrBadExtensions{extensions}
}
