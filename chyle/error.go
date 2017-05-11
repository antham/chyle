package chyle

import (
	"fmt"
)

// ErrWapper render a string error from an existing error
// with added string message
type ErrWapper struct {
	msg string
	err error
}

// Error dump error string
func (e ErrWapper) Error() string {
	return fmt.Sprintf("%s : %s", e.msg, e.err)
}

// addCustomMessageToError append an string message to an error
// by creating a brand new error
func addCustomMessageToError(msg string, err error) error {
	if err == nil {
		return nil
	}

	return ErrWapper{msg, err}
}
