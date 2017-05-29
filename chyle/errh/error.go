package errh

import (
	"fmt"
)

// errWrapper render a string error from an existing error
// with added string message
type errWrapper struct {
	msg string
	err error
}

// Error dump error string
func (e errWrapper) Error() string {
	return fmt.Sprintf("%s : %s", e.msg, e.err)
}

// AddCustomMessageToError append an string message to an error
// by creating a brand new error
func AddCustomMessageToError(msg string, err error) error {
	if err == nil {
		return nil
	}

	return errWrapper{msg, err}
}
