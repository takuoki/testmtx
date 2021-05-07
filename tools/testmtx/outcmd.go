package main

import (
	"errors"
	"fmt"

	"github.com/takuoki/gsheets/sheets"
	"github.com/takuoki/testmtx"
	"github.com/urfave/cli"
)

func init() {
	subCmdList = append(subCmdList, cli.Command{
		Name:  "out",
		Usage: "Outputs test data files",
		Action: func(c *cli.Context) error {
			return action(c, &output{})
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "sheet, s",
				Usage: "google spreadsheet id",
			},
			cli.StringFlag{
				Name:  "auth, a",
				Value: "credentials.json",
				Usage: "credential file for Google Sheets API",
			},
			cli.StringFlag{
				Name:  "xlsx, x",
				Usage: "excel file name",
			},
			cli.StringFlag{
				Name:  "format, f",
				Value: "json",
				Usage: "output format type (json, yaml)",
			},
			cli.StringFlag{
				Name:  "out, o",
				Value: "out",
				Usage: "output directory",
			},
			cli.IntFlag{
				Name:  "proplevel, pl",
				Value: 10,
				Usage: "property level (if you extend properties columns, mandatory)",
			},
			cli.StringFlag{
				Name:  "indent, i",
				Value: "  ",
				Usage: "indent string",
			},
		},
	})
}

type output struct{}

func (o *output) Run(c *cli.Context, conf *config) error {

	if c.String("sheet") == "" && c.String("xlsx") == "" {
		return errors.New("either the sheet option or the xlsx option is required")
	}
	if c.String("sheet") != "" && c.String("xlsx") != "" {
		return errors.New("specify only one of the 'sheet' option and the 'xlsx' option")
	}

	const inputTypeGss = "google-spreadsheet"
	const inputTypeExcel = "excel"
	inputType := ""
	if c.String("sheet") != "" {
		inputType = inputTypeGss
	}
	if c.String("xlsx") != "" {
		inputType = "excel"
	}

	p, err := testmtx.NewParser(testmtx.PropLevel(c.Int("proplevel")))
	if err != nil {
		return fmt.Errorf("unable to create parser: %w", err)
	}

	var f testmtx.Formatter
	switch c.String("format") {
	case "json":
		f, err = testmtx.NewJSONFormatter(testmtx.JSONIndentStr(c.String("indent")))
	case "yaml":
		f, err = testmtx.NewYamlFormatter(testmtx.YamlIndentStr(c.String("indent")))
	default:
		return fmt.Errorf("invalid format type (%s)", c.String("format"))
	}
	if err != nil {
		return fmt.Errorf("unable to create formatter: %w", err)
	}

	var d doc
	switch inputType {
	case inputTypeGss:
		sheetID := c.String("sheet")
		if v, ok := conf.SheetAliasMap[sheetID]; ok {
			sheetID = v
		}
		d, err = newGssDoc(sheetID, c.String("auth"))
		if err != nil {
			return fmt.Errorf("unable to create gss doc: %w", err)
		}
	case inputTypeExcel:
		d, err = newXlsxDoc(c.String("xlsx"))
		if err != nil {
			return fmt.Errorf("unable to create xlsx doc: %w", err)
		}
	}

	sheetNames, err := d.GetSheetNames()
	if err != nil {
		return fmt.Errorf("unable to retrieve sheet names: %w", err)
	}

	for _, sheetName := range sheetNames {
		if _, ok := conf.ExceptSheetSet[sheetName]; ok {
			continue
		}
		s, err := d.GetSheet(sheetName)
		if err != nil {
			return fmt.Errorf("unable to retrieve sheet data: %w", err)
		}
		sh, err := p.Parse(s, sheetName)
		if err != nil {
			return fmt.Errorf("unable to parse sheet data: %w", err)
		}
		err = testmtx.Output(f, sh, c.String("out"))
		if err != nil {
			return fmt.Errorf("unable to output test data: %w", err)
		}
	}

	fmt.Println("complete!")

	return nil
}

type doc interface {
	GetSheetNames() ([]string, error)
	GetSheet(sheetName string) (sheets.Sheet, error)
}
