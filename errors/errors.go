package errors

import (
	"fmt"
)

// UnexpectedError is used for unexpected errors
type UnexpectedError struct {
	Msg           string
	InternalError error
}

func (e *UnexpectedError) Error() string {
	return e.Msg
}

// NotFoundError happens when the document does not exist in the store
type NotFoundError struct {
	ID    string
	Model string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%v with id %q not found", e.Model, e.ID)
}

// BadRequestError happens when an id is not valid
type BadRequestError struct {
	Msg           string
	InternalError error
}

func (e *BadRequestError) Error() string {
	return e.Msg
}
