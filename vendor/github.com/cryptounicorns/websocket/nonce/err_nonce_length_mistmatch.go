package nonce

import (
	"fmt"
)

type ErrNonceLengthMismatch struct {
	Want int
	Got  int
}

func (e *ErrNonceLengthMismatch) Error() string {
	return fmt.Sprintf(
		"Nonce length mismatch, want '%d', got '%d'",
		e.Want,
		e.Got,
	)
}

func NewErrNonceLengthMismatch(w int, g int) *ErrNonceLengthMismatch {
	return &ErrNonceLengthMismatch{
		Want: w,
		Got:  g,
	}
}
