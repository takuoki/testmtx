package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/takuoki/gsheets"
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
				Usage: "google spreadsheet id (mandatory)",
			},
			cli.StringFlag{
				Name:  "auth, a",
				Value: "credentials.json",
				Usage: "credential file for Google Sheets API",
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

	if c.String("sheet") == "" {
		return errors.New("Please specify a google spreadsheet id")
	}

	p, err := testmtx.NewParser(testmtx.PropLevel(c.Int("proplevel")))
	if err != nil {
		return fmt.Errorf("Unable to create parser: %v", err)
	}

	var f testmtx.Formatter
	switch c.String("format") {
	case "json":
		f, err = testmtx.NewJSONFormatter(testmtx.JSONIndentStr(c.String("indent")))
	case "yaml":
		f, err = testmtx.NewYamlFormatter(testmtx.YamlIndentStr(c.String("indent")))
	default:
		return fmt.Errorf("Invalid format type (%s)", c.String("format"))
	}
	if err != nil {
		return fmt.Errorf("Unable to create formatter: %v", err)
	}

	ctx := context.Background()
	ctx = gsheets.WithCache(ctx)
	client, err := gsheets.NewForCLI(ctx, c.String("auth"))
	if err != nil {
		return fmt.Errorf("Unable to create gsheets client: %v", err)
	}

	sheetID := c.String("sheet")
	if v, ok := conf.SheetAliasMap[sheetID]; ok {
		sheetID = v
	}

	sheetNames, err := client.GetSheetNames(ctx, sheetID)
	if err != nil {
		return fmt.Errorf("Unable to retrieve sheet names: %v", err)
	}

	for _, sheetName := range sheetNames {
		if _, ok := conf.ExceptSheetSet[sheetName]; ok {
			continue
		}
		s, err := client.GetSheet(ctx, sheetID, sheetName)
		if err != nil {
			return fmt.Errorf("Unable to retrieve sheet data: %v", err)
		}
		sh, err := p.Parse(s, sheetName)
		if err != nil {
			return fmt.Errorf("Unable to parse sheet data: %v", err)
		}
		err = testmtx.Output(f, sh, c.String("out"))
		if err != nil {
			return fmt.Errorf("Unable to output test data: %v", err)
		}
	}

	fmt.Println("complete!")

	return nil
}
