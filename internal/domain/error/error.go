package domainerror

import (
	"errors"
	"fmt"
)

type Kind string

const (
	KindValidation Kind = "validation_error"
	KindNotFound   Kind = "not_found"
	KindConflict   Kind = "conflict"
	KindInternal   Kind = "internal_error"
)

type Error struct {
	Kind    Kind
	Message string
	Fields  map[string]string
	Err     error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return string(e.Kind)
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func cloneFields(fields map[string]string) map[string]string {
	if len(fields) == 0 {
		return nil
	}
	out := make(map[string]string, len(fields))
	for k, v := range fields {
		out[k] = v
	}
	return out
}

func Validation(message string, fields map[string]string) *Error {
	if message == "" {
		message = "validation failed"
	}
	return &Error{
		Kind:    KindValidation,
		Message: message,
		Fields:  cloneFields(fields),
	}
}

func NotFound(resource string) *Error {
	message := "resource not found"
	if resource != "" {
		message = fmt.Sprintf("%s not found", resource)
	}
	return &Error{
		Kind:    KindNotFound,
		Message: message,
	}
}

func Conflict(message string) *Error {
	if message == "" {
		message = "conflict"
	}
	return &Error{
		Kind:    KindConflict,
		Message: message,
	}
}

func Internal(err error) *Error {
	if err == nil {
		err = errors.New("internal error")
	}
	return &Error{
		Kind:    KindInternal,
		Message: "internal server error",
		Err:     err,
	}
}

func As(err error) *Error {
	if err == nil {
		return nil
	}
	var derr *Error
	if errors.As(err, &derr) {
		return derr
	}
	return Internal(err)
}

