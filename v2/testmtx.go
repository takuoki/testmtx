// Package testmtx helps you to create test data files with Google Spreadsheets.
// Once you create test cases as matrix on Google Spreadsheets,
// this tool generates test data like JSON based on the data you input.
// Using testmtx, you can get advantages of `completeness`, `readability` and `consistency` for testing.
//
// This package is just liblary.
// If what you want is just to use, see the standard tool below which uses this package.
// - github.com/takuoki/testmtx/tools/testmtx
package testmtx

// Property type in spreadsheet.
const (
	typeObject = "object"
	typeArray  = "array"
	typeString = "string"
	typeNumber = "number"
	typeBool   = "bool"
)

// String for special value.
// TODO: Optionにするか
const (
	strNull  = "*null"
	strNew   = "*new"
	strEmpty = "*empty"
)

// これはParserでいいな
// type Client struct {
// 	convertSimpleValueFuncs map[string]func(s string) (SimpleValue, error)
// }

// func New() *Client {
// 	return &Client{
// 		convertSimpleValueFuncs: defaultConvertSimpleValueFuncs,
// 	}
// }

// Client Option
