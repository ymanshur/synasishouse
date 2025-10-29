package typex

import (
	"fmt"

	"github.com/ymanshur/synasishouse/pkg/util"
)

// UnProcessableEnity status error
type UnProcessableEnity string

func (s UnProcessableEnity) Error() string {
	return string(s)
}

// Conflict status error
type Conflict string

func (s Conflict) Error() string {
	return string(s)
}

// NotFound status error
type NotFound string

func (s NotFound) Error() string {
	return string(s)
}

// Unavailable status error
type Unavailable string

func (s Unavailable) Error() string {
	return string(s)
}

// NewUnprocessableEntityError creates [UnProcessableEnity] status with custom error message
func NewUnprocessableEntityError(s string) error {
	return UnProcessableEnity(s)
}

// NewConflictError creates [Conflict] status with custom error message
func NewConflictError(data string) error {
	return Conflict(fmt.Sprintf("%s already exists", data))
}

// NewNotFoundError creates [NotFound] status with custom error message and wrap predecessor errors
func NewNotFoundError(data string, errs ...error) error {
	s := NotFound(fmt.Sprintf("%s not found", data))
	return util.JoinErrors(s, errs...)
}

// NewUnavailableError creates [Unavailable] status with custom error message
func NewUnavailableError(service string, errs ...error) error {
	s := Unavailable(fmt.Sprintf("%s unavailable", service))
	return util.JoinErrors(s, errs...)
}
