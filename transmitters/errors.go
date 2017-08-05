package transmitters

import (
	"fmt"
)

type ErrUnknownTransmitterType struct {
	t string
}

func (e *ErrUnknownTransmitterType) Error() string {
	return fmt.Sprintf(
		"Unknown transmitter type '%s'",
		e.t,
	)
}
func NewErrUnknownTransmitterType(t string) error {
	return &ErrUnknownTransmitterType{t}
}

//
