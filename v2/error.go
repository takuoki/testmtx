package testmtx

import (
	"fmt"

	"github.com/takuoki/clmconv"
)

// ParseError is an error type that occurs during parsing.
// `msg` is required, and the others are optional.
// Set where you can set it, and wrap it.
type ParseError struct {
	msg       string
	sheet     *string
	rowNumber *int
	clmLetter *string
	err       error
}

// newParseError returns a new ParseError.
func newParseError(msg string, options ...parseErrorOption) *ParseError {
	e := &ParseError{msg: msg}
	for _, opt := range options {
		opt(e)
	}
	return e
}

// Error returns the error message.
// If the error is wrapped, the original error message is returned.
// The error message includes the sheet name and the cell position.
// If the error is wrapped, the sheet name and the cell position
// as close as possible to the original error are included.
func (e *ParseError) Error() string {

	msg := ""
	sheet := ""
	row := 0
	column := ""

	var err error = e
	for {
		if x, ok := err.(*ParseError); ok {
			msg = x.msg
			if x.sheet != nil {
				sheet = *x.sheet
			}
			if x.rowNumber != nil {
				row = *x.rowNumber
			}
			if x.clmLetter != nil {
				column = *x.clmLetter
			}
		} else {
			msg = err.Error()
		}

		if x, ok := err.(interface{ Unwrap() error }); ok {
			err = x.Unwrap()
			if err == nil {
				break
			}
		} else {
			break
		}
	}

	return fmt.Sprintf("%s (sheet=%q, cell=\"%s%d\")", msg, sheet, column, row)
}

// DetailError returns the detail error message.
// If the error is wrapped, all error messages are concatenated.
// The error message includes the sheet name and the cell position.
// If the error is wrapped, the sheet name and the cell position
// as close as possible to the original error are included.
func (e *ParseError) DetailError() string {

	msg := ""
	sheet := ""
	row := 0
	column := ""

	var err error = e
	for {
		if x, ok := err.(*ParseError); ok {
			if msg == "" {
				msg = x.msg
			} else {
				msg += fmt.Sprintf(": %s", x.msg)
			}
			if x.sheet != nil {
				sheet = *x.sheet
			}
			if x.rowNumber != nil {
				row = *x.rowNumber
			}
			if x.clmLetter != nil {
				column = *x.clmLetter
			}
		} else {
			msg += fmt.Sprintf(": %s", err.Error())
		}

		if x, ok := err.(interface{ Unwrap() error }); ok {
			err = x.Unwrap()
			if err == nil {
				break
			}
		} else {
			break
		}
	}

	return fmt.Sprintf("%s (sheet=%q, cell=\"%s%d\")", msg, sheet, column, row)
}

func (e *ParseError) Unwrap() error {
	return e.err
}

// parseErrorOption changes some parameters of the ParseError.
type parseErrorOption func(*ParseError)

// parseErrorOptionSheetName changes the sheet name in the ParseError.
func parseErrorOptionSheetName(name string) parseErrorOption {
	return func(e *ParseError) {
		e.sheet = &name
	}
}

// parseErrorOptionRowNumber changes the row number in the ParseError.
// The row number is 1-based.
func parseErrorOptionRowNumber(number int) parseErrorOption {
	return func(e *ParseError) {
		e.rowNumber = &number
	}
}

// parseErrorOptionColumnLetterIndex changes the column letter in the ParseError.
// The column letter index is 0-based.
func parseErrorOptionColumnLetterIndex(index int) parseErrorOption {
	return func(e *ParseError) {
		letter := clmconv.Itoa(index)
		e.clmLetter = &letter
	}
}

// parseErrorOptionBaseError changes the base error in the ParseError.
func parseErrorOptionBaseError(err error) parseErrorOption {
	return func(e *ParseError) {
		e.err = err
	}
}
