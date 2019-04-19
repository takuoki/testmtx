# testmtx

A golang package and tool for converting test cases written on Google Spreadsheets to test data files like JSON.

<!-- vscode-markdown-toc -->
* [Description](#Description)
* [Installation](#Installation)
* [Requirements](#Requirements)
* [Usage](#Usage)
	* [Output Test Data Files](#OutputTestDataFiles)
		* [1. Create test cases](#Createtestcases)
		* [2. Execute command](#Executecommand)
		* [3. Check generated files](#Checkgeneratedfiles)
	* [Output Property](#OutputProperty)
		* [1. Create Go type](#CreateGotype)
		* [2. Execute command](#Executecommand-1)
* [Config File](#ConfigFile)

<!-- vscode-markdown-toc-config
	numbering=false
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

## <a name='Description'></a>Description

`testmtx` helps you to create test data files with Google Spreadsheets.
Once you create test cases as matrix on Google Spreadsheets, this tool generates test data like JSON based on the data you input.
Using `testmtx`, you can get advantages of **completeness**, **readability** and **consistency** for testing.

## <a name='Installation'></a>Installation

If you have a Go environment, you can install this tool using `go get`. Before installing, enable the Go module feature.

```bash
go get github.com/takuoki/testmtx/tools/testmtx
```

If not, download it from [the release page](https://github.com/takuoki/testmtx/releases).

## <a name='Requirements'></a>Requirements

This tool uses Google OAuth2.0. So before executing tool, you have to prepare `credentials.json`. See [Go Quickstart](https://developers.google.com/sheets/api/quickstart/go), or [Blog (Japanese)](https://medium.com/veltra-engineering/how-to-use-google-sheets-api-with-golang-9e50ee9e0abc) for the details.

This is brief steps.

  1. Create new GCP Project on [GCP Console](https://console.cloud.google.com)
  1. Enable Google Sheets API on [APIs & Services] - [API Library]
  1. Setting OAuth consent screen and create OAuth client ID on [APIs & Services] - [Credentials]
  1. Download JSON file and rename it to `credentials.json`

## <a name='Usage'></a>Usage

### <a name='OutputTestDataFiles'></a>Output Test Data Files

#### <a name='Createtestcases'></a>1. Create test cases

Copy [the Sample Sheet](https://docs.google.com/spreadsheets/d/1Zs2HI7x8eQ05ICoaBdv1I1ny_KtmtrE05Lyb7OwYmdE) and fill it as you want.

* Property Area

  The left side of the sheet is the property area.
  You should input all property definitions here, and if the property is `array` type, repeat it as you need.
  Root property names is not used in the output files. This is used as the output directory name.
  If you already have a Go type for your test data, you can generate a property list using [`prop` sub command](#OutputProperty).
  If the property level is not enough at the default 10, add columns. In this case, adjust at run time with the `-proplevel` or `-pl` option.

* Case Area

  The right side of the sheet is the case area. Each column is one test case.
  You should input case names and test data for that case.
  A blank cells mean `null`, so the property itself is not output.
  If you want to specify `null` clearly, you can also use `*null` keyword.
  When you want to output objects or arrays, and empty strings, use `*new` and `*empty` keywords.

![Sample Sheet](https://github.com/takuoki/testmtx/blob/image/image/sample_sheet.png)

#### <a name='Executecommand'></a>2. Execute command

Using `out` sub command, you can generate test data with Google Spreadsheets.
This tool creates test data for all sheets, and all test cases.
If you want to ignore some sheets, use the except sheet name feature in configuration.
For each sheet, this tool searches from the beginning of the test case name to the right and end when the test case name becomes blank.

```bash
$ testmtx -c config.json out -s sample
complete!
```

sample : sheet ID alias for `1Zs2HI7x8eQ05ICoaBdv1I1ny_KtmtrE05Lyb7OwYmdE` (see [Config File](#config-file))

#### <a name='Checkgeneratedfiles'></a>3. Check generated files

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

### <a name='OutputProperty'></a>Output Property

#### <a name='CreateGotype'></a>1. Create Go type

If you already have a Go type for your test data like below, you can generate a property list with this tool.

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

#### <a name='Executecommand-1'></a>2. Execute command

Using `prop` sub command, you can generate a property list for Google Spreadsheets from Go type.
This is a subsidiary function, and you can modify this output list.
If the target type refers some other files, you should modify the output type.

```bash
$ testmtx prop -f sample/sample.go -t Request
request             object
    num_key         numder
    string_key      string
    bool_key        bool
    object_key      object
        key1        numder
        key2        string
    array_key       array
        * 0         object
            key3    numder
            key4    string
        * 1         object
            key3    numder
            key4    string

complete!
```

## <a name='ConfigFile'></a>Config File

You can use several additional functions with the config file.
When you want to use these functions, specify config file as command line argument.

```json
{
  "except_sheet_names": [
    "overview"
  ],
  "sheet_list": [
    {
      "name": "testmtx_sample",
      "alias": "sample",
      "sheet_id": "1Zs2HI7x8eQ05ICoaBdv1I1ny_KtmtrE05Lyb7OwYmdE"
    }
  ]
}
```

* except_sheet_names: The sheets listed here are excluded from output.
* sheet_list: If you define an alias here, you can specify a sheet with an alias name.

You can check the contents of configuration using `conf` sub command.

```bash
$ testmtx -c config.json conf
# Except Sheet Names
- overview

# Sheet List
  NAME           | ALIAS  | SPREADSHEET ID
--------------------------------------------------------------------------
  testmtx_sample | sample | 1Zs2HI7x8eQ05ICoaBdv1I1ny_KtmtrE05Lyb7OwYmdE
```
