# testmtx

`testmtx` is a test data generator using Google Spreadsheets.

## Description

`testmtx` helps you to create test data using Google Spreadsheets.
If you create test case as matrix on Google Spreadsheets, this tool generates some test data like JSON file based on the data you input.
Using `testmtx`, you can get advantages of **completeness**, **readability** and **consistency**.

## Install

TBD...

## Preparation

* credentials.json for Google Sheets API
* config.json for `testmtx` (TBD)

## Usage

### Google Spreadsheets

image...

### Command

```txt
./testmtx out -s {SpreadsheetID}
```

### Generated File

./out/request/sheetname_casename.json

```json
{
  "val": 111,
  "str": "string value 1",
  "b": true,
  "obj": {
    "foo": 222,
    "bar": "string value 2"
  },
  "ary": [
    {
      "foo": 333,
      "bar": "string value 3"
    },
    {
      "foo": 444,
      "bar": "string value 4"
    }
  ]
}
```

./out/expected/sheetname_casename.json

```json
{
  "code": "CODE",
  "message": "MESSAGE"
}
```

## How to input sheet

* If object or array type, use `*new` keyword
* If string type and empty, use `*empty` keyword

## Other Requirements

* create test data until case name become empty
* empty cell is nil
