package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/takuoki/testmtx/v2"
	"github.com/takuoki/testmtx/v2/cmd/helper"
	cli "github.com/urfave/cli/v2"
)

const (
	appName = "testmtx"
	version = "2.0.0"
)

func main() {
	outCommand, err := outCommandFunc()
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:    appName,
		Version: version,
		Usage:   "A test data generator using spreadsheet.",
		Commands: []*cli.Command{
			outCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// TODO: 全体の構築自体をhelperに任せたい
var outCommandFunc = func() (*cli.Command, error) {

	parseFlags, newParserFunc, err := helper.GetParseFlagAndFunc(nil)
	if err != nil {
		return nil, fmt.Errorf("fail to get parse flag and func: %w", err)
	}

	// TODO: indentStrはデフォルトありにしたい
	formatFlags, newFormatterFunc, err := helper.GetFormatFlagAndFunc([]helper.Formatter{
		{Name: "json", NewFunc: testmtx.NewJSONFormatter},
		{Name: "yaml", NewFunc: testmtx.NewYAMLFormatter},
	}, "  ")
	if err != nil {
		return nil, fmt.Errorf("fail to get format flag and func: %w", err)
	}

	layoutFlags, newOutputterFunc, err := helper.GetLayoutFlagAndFunc([]helper.Outputter{
		{Name: "1column-1case", NewFunc: testmtx.NewOneColumnOneCaseOutputter},
		{Name: "1sheet-1case", NewFunc: testmtx.NewOneSheetOneCaseOutputter},
	})
	if err != nil {
		return nil, fmt.Errorf("fail to get layout flag and func: %w", err)
	}

	// TODO: flagの順をきれいにしたい
	var flags []cli.Flag
	flags = append(flags, parseFlags...)
	flags = append(flags, formatFlags...)
	flags = append(flags, layoutFlags...)
	flags = append(flags,
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Value:   "excel",
			Usage:   `input spreadsheet type ("excel" or "gs")`,
		},
		&cli.StringFlag{
			Name:     "filepath",
			Aliases:  []string{"f"},
			Required: true,
			Usage:    `input spreadsheet filepath (type="excel"), or google spreadsheet ID (type="gs")`,
		},
		&cli.StringFlag{
			Name:     "sheet",
			Aliases:  []string{"s"},
			Required: true,
			Usage:    "input sheet name",
		},
		&cli.StringFlag{
			Name:    "out",
			Aliases: []string{"o"},
			Value:   "out",
			Usage:   "output root directory",
		},
	)

	return &cli.Command{
		Name:  "out",
		Usage: "Outputs test data files",
		Flags: flags,
		Action: func(c *cli.Context) error {

			parser, err := newParserFunc(c)
			if err != nil {
				return fmt.Errorf("fail to create parser: %w", err)
			}

			formatter, err := newFormatterFunc(c)
			if err != nil {
				return fmt.Errorf("fail to create formatter: %w", err)
			}

			outputter, err := newOutputterFunc(c, formatter)
			if err != nil {
				return fmt.Errorf("fail to create outputter: %w", err)
			}

			if c.String("type") != "excel" {
				return errors.New("unsupportted type")
			}

			doc, err := testmtx.NewExcelDoc(c.String("filepath"))
			if err != nil {
				return fmt.Errorf("fail to create xlsx doc: %w", err)
			}

			docSheet, err := doc.GetSheet(c.String("sheet"))
			if err != nil {
				return fmt.Errorf("fail to get sheet: %w", err)
			}

			sheet, err := parser.Parse(docSheet)
			if err != nil {
				return fmt.Errorf("fail to parse sheet: %w", err)
			}

			if err := outputter.Output(c.String("out"), sheet); err != nil {
				return fmt.Errorf("fail to output: %w", err)
			}

			return nil
		},
	}, nil
}
