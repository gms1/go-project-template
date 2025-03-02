package core

import (
	"errors"
	"fmt"
)

var ErrSignalHandlerAlreadyInitialized = errors.New("signal handler is already initialized")

func ErrorRootCause(err error) error {
	result := err
	for cause := result; cause != nil; cause = errors.Unwrap(cause) {
		result = cause
	}
	return result
}

type stackTraceError struct {
	text     string
	stack    []byte
	causedBy error
}

func (e *stackTraceError) Error() string {
	return e.text
}

func (e *stackTraceError) Stack() string {
	return string(e.stack)
}

func (e *stackTraceError) Unwrap() error {
	return e.causedBy
}

func ToStackTraceError(stack []byte, err error) error {
	return &stackTraceError{err.Error(), stack, errors.Unwrap(err)}
}

func StackTraceErrorf(stack []byte, s string, vals ...any) error {
	return ToStackTraceError(stack, fmt.Errorf(s, vals...))
}

func NewStackTraceError(stack []byte, v any) error {
	return StackTraceErrorf(stack, "%v", v)
}

func Stack(err error) error {
	u, ok := err.(interface {
		Stack() string
	})
	if !ok {
		return nil
	}
	return errors.New(u.Stack())
}
