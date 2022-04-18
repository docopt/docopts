package docopt

import (
	"fmt"
)

type ErrorType int

const (
	ErrorUser ErrorType = iota
	ErrorLanguage
)

func (e ErrorType) String() string {
	switch e {
	case ErrorUser:
		return "ErrorUser"
	case ErrorLanguage:
		return "ErrorLanguage"
	}
	return ""
}

// UserError records an error with program arguments.
type UserError struct {
	msg   string
	Usage string
}

func (e UserError) Error() string {
	return e.msg
}
func newUserError(msg string, f ...interface{}) error {
	return &UserError{fmt.Sprintf(msg, f...), ""}
}

// LanguageError records an error with the doc string.
type LanguageError struct {
	msg string
}

func (e LanguageError) Error() string {
	return e.msg
}
func newLanguageError(msg string, f ...interface{}) error {
	return &LanguageError{fmt.Sprintf(msg, f...)}
}

var newError = fmt.Errorf
