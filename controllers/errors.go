package controllers

import (
	"fmt"
)

// InvalidBodyError occurs when the body's request is invalid
type InvalidBodyError struct {
	InternalError error
}

func (e *InvalidBodyError) Error() string {
	return fmt.Sprint("Invalid body")
}

// NoBodyError occurs when the request does not have body
type NoBodyError struct{}

func (e *NoBodyError) Error() string {
	return fmt.Sprint("No body")
}
