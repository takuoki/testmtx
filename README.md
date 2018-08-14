# testmtx

`testmtx` is a test data generator using Google Spreadsheets.

## Description

`testmtx` helps you to create test data using Google Spreadsheets.
If you create test case as matrix on Google Spreadsheets, this tool generates some test data like JSON file based on the data you input.
Using `testmtx`, you can get advantages of **completeness**, **readability** and **consistency**.

## Installation

To install `testmtx`, please use `go get`.
You must have Go 1.9 or greater installed, and `$GOPATH/bin` added to your PATH.

```txt
$ go get github.com/takuoki/testmtx
...
$ testmtx -h
...
```

## Preparation

* `credentials.json` for [Google Sheets API](https://developers.google.com/sheets/api/quickstart/go#step_1_turn_on_the)

  1. Create new GCP Project on [GCP Console](https://console.cloud.google.com)
  1. Enable Google Sheets API on [APIs & Services] - [API Library]
  1. Setting OAuth consent screen and create OAuth client ID on [APIs & Services] - [Credentials]
  1. Download JSON file and rename it to `credentials.json`

## Usage

### Output Test Data Files

Using `out` subcommand, you can generate test data from Google Spreadsheets.
This tool creates test data for all sheets, and all test cases.
For each sheet, this tool searches from the beginning of the test case name to the right and end when the test case name becomes blank.

Blank cells mean `null`, so the property itself is not output.
When you want to output object or array, or empty characters in string, use `*new`,　`*empty` keywords.

#### Google Spreadsheets

[Sample Sheet](https://docs.google.com/spreadsheets/d/1Zs2HI7x8eQ05ICoaBdv1I1ny_KtmtrE05Lyb7OwYmdE)

![Sample Sheet](https://github.com/takuoki/testmtx/blob/image/image/sample_sheet.png)

#### Command: out

```txt
$ testmtx -c config.json out -s sample
output completed successfully!
```

sample : sheet ID alias for `1Zs2HI7x8eQ05ICoaBdv1I1ny_KtmtrE05Lyb7OwYmdE` (see [Config File](#config-file))

#### Generated File

./out/request/sheet_case1.json

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

./out/expected/sheet_case1.json

```json
{
  "status": "success",
  "code": 200
}
```

### Output Property

Using `prop` subcommand, you can generate the list of properties for Google Spreadsheets from Go type.
This is a subsidiary function, and you can modify this output list.
If the target type refers some other files, you should modify the output type.

#### Golang File

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
$ testmtx prop -f sample/sample.go -t Request
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

output completed successfully!
```

## Config File

You can use several additional functions with the config file.
When you want to use these functions, specify config file as command line argument.

```json
{
  "except_sheet_names": [
    "overview"
  ],
  "sheet_alias_list": [
    {
      "alias": "sample",
      "sheet_id": "1Zs2HI7x8eQ05ICoaBdv1I1ny_KtmtrE05Lyb7OwYmdE"
    }
  ]
}
```

* except_sheet_names: The sheets listed here are excluded from output.
* sheet_alias_list: If you define an alias here, you can specify a sheet with an alias name.
