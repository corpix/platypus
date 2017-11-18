package errors

import (
	"fmt"
)

type ErrUnknownRouterType struct {
	t string
}

func (e *ErrUnknownRouterType) Error() string {
	return fmt.Sprintf(
		"Unknown router type '%s'",
		e.t,
	)
}
func NewErrUnknownRouterType(t string) error {
	return &ErrUnknownRouterType{t}
}
