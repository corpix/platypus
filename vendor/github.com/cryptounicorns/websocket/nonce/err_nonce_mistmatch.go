package nonce

import (
	"fmt"
)

type ErrNonceMismatch struct {
	Want Nonce
	Got  Nonce
}

func (e *ErrNonceMismatch) Error() string {
	return fmt.Sprintf(
		"Nonce mismatch, want '%s', got '%s'",
		e.Want,
		e.Got,
	)
}

func NewErrNonceMismatch(w Nonce, g Nonce) *ErrNonceMismatch {
	return &ErrNonceMismatch{
		Want: w,
		Got:  g,
	}
}
