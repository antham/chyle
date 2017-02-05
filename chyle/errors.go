package chyle

import (
	"fmt"
)

// ErrTemplateMalformed issued when something wrong
// happened when creation a template
type ErrTemplateMalformed struct {
	err error
}

// Error return string error
func (e ErrTemplateMalformed) Error() string {
	return fmt.Sprintf("check your template is well-formed : %s", e.err.Error())
}

// ErrSenderFailed issued when something went
// wrong during third part service request
type ErrSenderFailed struct {
	reason string
}

// Error return string error
func (e ErrSenderFailed) Error() string {
	return fmt.Sprintf("sender issue : %s", e.reason)
}
