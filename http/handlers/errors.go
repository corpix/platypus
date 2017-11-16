package handlers

import (
	"fmt"
)

type ErrUnknownHandlerType struct {
	t string
}

func (e *ErrUnknownHandlerType) Error() string {
	return fmt.Sprintf(
		"Unknown handler type '%s'",
		e.t,
	)
}
func NewErrUnknownHandlerType(t string) error {
	return &ErrUnknownHandlerType{t}
}
