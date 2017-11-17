package errors

import (
	"fmt"
)

type ErrBadSubProtocol struct {
	SubProtocol string
}

func (e *ErrBadSubProtocol) Error() string {
	return fmt.Printf(
		"Bad sub-protocol: %s",
		e.SubProtocol,
	)
}

func NewErrBadSubProtocol(subProtocol string) *ErrBadSubProtocol {
	return &ErrBadSubProtocol{subProtocol}
}
