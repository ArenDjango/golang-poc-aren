package helper

import "fmt"

const (
	Unknown            = 2
	DeadlineExceeded   = 4
	NotFound           = 5
	AlreadyExists      = 6
	PermissionDenied   = 7
	ResourceExhausted  = 8
	FailedPrecondition = 9
	InternalError      = 10
)

type BusinessError struct {
	Status uint8
	Err    error
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("status: %d, error: %v", e.Status, e.Err)
}

func (e *BusinessError) Unwrap() error {
	return e.Err
}

func NewError(status uint8, err error) *BusinessError {
	return &BusinessError{
		Status: status,
		Err:    err,
	}
}
