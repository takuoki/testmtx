package helper

import (
	"errors"
	"fmt"

	cli "github.com/urfave/cli/v2"

	"github.com/takuoki/testmtx/v2"
)

func GetOutCommand(options ...GetOutCommandOption) (*cli.Command, error) {
	opts := defaultGetOutCommandOptions
	for _, o := range options {
		o(&opts)
	}

	parseFlags, newParserFunc, err := getParseFlagAndFunc(opts.additionalSimpleValues, opts.defaultPropLevel)
	if err != nil {
		return nil, fmt.Errorf("fail to get parse flag and func: %w", err)
	}

	formatFlags, newFormatterFunc, err := getFormatFlagAndFunc(opts.formatters, opts.indentStr)
	if err != nil {
		return nil, fmt.Errorf("fail to get format flag and func: %w", err)
	}

	layoutFlags, newOutputterFunc, err := getLayoutFlagAndFunc(opts.layouts)
	if err != nil {
		return nil, fmt.Errorf("fail to get layout flag and func: %w", err)
	}

	flags := []cli.Flag{
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
	}
	flags = append(flags, parseFlags...)
	flags = append(flags, formatFlags...)
	flags = append(flags, layoutFlags...)

	return &cli.Command{
		Name:  "out",
		Usage: "Outputs test data files",
		Flags: flags,
		Action: func(c *cli.Context) error {

			if c.String("type") != "excel" {
				return errors.New("unsupportted type")
			}

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

type getOutCommandOptions struct {
	additionalSimpleValues map[string]testmtx.ConvertValueFunc
	defaultPropLevel       int
	formatters             []Formatter
	indentStr              string
	layouts                []Layout
}

var defaultGetOutCommandOptions = getOutCommandOptions{
	additionalSimpleValues: nil,
	defaultPropLevel:       10,
	formatters: []Formatter{
		{Name: "json", NewFunc: testmtx.NewJSONFormatter},
		{Name: "yaml", NewFunc: testmtx.NewYAMLFormatter},
	},
	indentStr: "  ",
	layouts: []Layout{
		{Name: "1column-1case", NewFunc: testmtx.NewOneColumnOneCaseOutputter},
		{Name: "1sheet-1case", NewFunc: testmtx.NewOneSheetOneCaseOutputter},
	},
}

type GetOutCommandOption func(*getOutCommandOptions)

func AdditionalSimpleValues(convertValueFuncs map[string]testmtx.ConvertValueFunc) GetOutCommandOption {
	return func(o *getOutCommandOptions) {
		o.additionalSimpleValues = convertValueFuncs
	}
}

func DefaultPropLevel(level int) GetOutCommandOption {
	return func(o *getOutCommandOptions) {
		o.defaultPropLevel = level
	}
}

func Formatters(formatters []Formatter) GetOutCommandOption {
	return func(o *getOutCommandOptions) {
		o.formatters = formatters
	}
}

func IndentStr(indentStr string) GetOutCommandOption {
	return func(o *getOutCommandOptions) {
		o.indentStr = indentStr
	}
}

func Layouts(layouts []Layout) GetOutCommandOption {
	return func(o *getOutCommandOptions) {
		o.layouts = layouts
	}
}
