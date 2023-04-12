package testmtx

import "fmt"

// TODO: コメント
type ParseError struct {
	msg       string
	sheet     *string
	rowNumber *int    // 1始まりの番号
	clmLetter *string // アルファベット
	err       error
}

// TODO: コメント
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

// TODO: コメント
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
				msg = fmt.Sprintf("%s: %s", msg, x.msg)
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
			msg = fmt.Sprintf("%s: %s", msg, err.Error())
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

// TODO: これをどこに配置するか
func pointer[T any](value T) *T {
	return &value
}
