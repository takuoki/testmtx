# testmtx

`testmtx` is a test data generator using Google Spreadsheets.

## Description

`testmtx` helps you to create test data using Google Spreadsheets.
If you create test case as matrix on Google Spreadsheets, this tool generates some test data like JSON file based on the data you input.
Using `testmtx`, you can get advantages of **completeness**, **readability** and **consistency**.

## Install

TBD...

## Preparation

* `credentials.json` for [Google Sheets API](https://developers.google.com/sheets/api/quickstart/go#step_1_turn_on_the)
* `config.json` for `testmtx` (TBD)

## Usage

### Output Property

#### Struct File

```go
type Request struct {
  NumKey    *int    `json:"num_key"`
  StringKey *string `json:"string_key"`
  BoolKey   *bool   `json:"bool_key"`
  ObjectKey struct {
    Key1 *int    `json:"key1"`
    Key2 *string `json:"key2"`
  } `json:"object_key"`
  ArrayKey []struct {
    Key3 *int    `json:"key3"`
    Key4 *string `json:"key4"`
  } `json:"array_key"`
}
```

#### Command: prop

```txt
./testmtx prop -f sample/sample.go -s Request
```

#### Output Format

Output to standard output.

```txt
request             object
    num_key         numder
    string_key      string
    bool_key        bool
    object_key      object
        key1        numder
        key2        string
    array_key       array
        *           object
            key3    numder
            key4    string
        *           object
            key3    numder
            key4    string
```

### Output Test Data Files

#### Google Spreadsheets

[Sample Sheet](https://docs.google.com/spreadsheets/d/1Zs2HI7x8eQ05ICoaBdv1I1ny_KtmtrE05Lyb7OwYmdE)

![Sample Sheet](https://github.com/takuoki/testmtx/blob/image/image/sample_sheet.png)

#### Command: out

```txt
./testmtx out -s {SpreadsheetID}
```

#### Generated File

./out/request/sheetname_casename.json

```json
{
  "num_key": 101,
  "string_key": "string value 101",
  "bool_key": true,
  "object_key": {
    "key1": 201,
    "key2": "string value 201"
  },
  "array_key": [
    {
      "key3": 301,
      "key4": "string value 301"
    },
    {
      "key3": 401,
      "key4": "string value 401"
    }
  ]
}
```

./out/expected/sheetname_casename.json

```json
{
  "status": "success",
  "code": 200
}
```

## How to input sheet

* If object or array type, use `*new` keyword
* If string type and empty, use `*empty` keyword

## Other Requirements

* create test data until case name become empty
* empty cell is nil
