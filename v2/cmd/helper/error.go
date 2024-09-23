package helper

import (
	"errors"
)

// UserError is an error caused by the user who executed the tool.
// It is not an error caused by the way the tool was made.
type UserError struct {
	msg string
}

func (e *UserError) Error() string {
	return e.msg
}

func ExtractErrorMsg(err error) string {
	if uerr, ok := extractUserError(err); ok {
		return uerr.Error()
	}
	return "!!!!! THIS MAY BE AN ERROR RELATED TO THE WAY THE TOOL WAS MADE !!!!!\n" +
		err.Error()
}

// extractUserError extracts the UserError from the error.
func extractUserError(err error) (userError error, ok bool) {
	var uerr *UserError
	if errors.As(err, &uerr) {
		return uerr, true
	}
	return nil, false
}
