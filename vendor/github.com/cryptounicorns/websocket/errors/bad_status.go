package errors

import (
	"fmt"
)

type ErrBadStatus struct {
	Status int
}

func (e *ErrBadStatus) Error() string {
	return fmt.Printf(
		"Bad HTTP status: %d",
		e.Status,
	)
}

func NewErrBadStatus(status int) *ErrBadStatus {
	return &ErrBadStatus{status}
}
