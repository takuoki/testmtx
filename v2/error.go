package testmtx

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/takuoki/clmconv"
)

// OpenFileError is an error type that occurs when opening a file.
type OpenFileError struct {
	filepath string
	err      error
}

func (e *OpenFileError) Error() string {
	var serr syscall.Errno
	if errors.As(e.err, &serr) {
		return fmt.Sprintf("fail to open file (path=%q): %s", e.filepath, serr)
	}
	return fmt.Sprintf("fail to open file (path=%q): %s", e.filepath, e.err)
}

func (e *OpenFileError) Unwrap() error {
	return e.err
}

// CreateFileError is an error type that occurs when creating a file.
type CreateFileError struct {
	filepath string
	err      error
}

func (e *CreateFileError) Error() string {
	var serr syscall.Errno
	if errors.As(e.err, &serr) {
		return fmt.Sprintf("fail to create file (path=%q): %s", e.filepath, serr)
	}
	return fmt.Sprintf("fail to create file (path=%q): %s", e.filepath, e.err)
}

func (e *CreateFileError) Unwrap() error {
	return e.err
}

// CreateDirError is an error type that occurs when creating a directory.
type CreateDirError struct {
	dirpath string
	err     error
}

func (e *CreateDirError) Error() string {
	var serr syscall.Errno
	if errors.As(e.err, &serr) {
		return fmt.Sprintf("fail to create directory (path=%q): %s", e.dirpath, serr)
	}
	return fmt.Sprintf("fail to create directory (path=%q): %s", e.dirpath, e.err)
}

func (e *CreateDirError) Unwrap() error {
	return e.err
}

// NotFoundError is an error type that occurs when the target is not found.
type NotFoundError struct {
	msg string
}

func (e *NotFoundError) Error() string {
	return e.msg
}

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
