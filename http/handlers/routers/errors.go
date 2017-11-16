package routers

import (
	"fmt"
)

type ErrUnknownTransmitterType struct {
	t string
}

func (e *ErrUnknownTransmitterType) Error() string {
	return fmt.Sprintf(
		"Unknown router type '%s'",
		e.t,
	)
}
func NewErrUnknownTransmitterType(t string) error {
	return &ErrUnknownTransmitterType{t}
}

//
