package helper

import (
	"errors"
	"fmt"

	cli "github.com/urfave/cli/v2"

	"github.com/takuoki/testmtx/v2"
)

func NewOutCommand(options ...NewOutCommandOption) (*cli.Command, error) {
	opts := defaultNewOutCommandOptions
	for _, o := range options {
		o(&opts)
	}

	getSheetNamesFlags, getSheetNamesFunc := opts.getSheetNamesFlagAndFunc()

	parseFlags, newParserFunc, err := getParseFlagAndFunc(opts.additionalSimpleValues, opts.defaultPropLevel)
	if err != nil {
		return nil, fmt.Errorf("fail to get parse flag and func: %w", err)
	}

	formatFlags, newFormatterFunc, err := getFormatFlagAndFunc(formatters, opts.indentStr)
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
	}
	flags = append(flags, getSheetNamesFlags...)
	flags = append(flags, parseFlags...)
	flags = append(flags, formatFlags...)
	flags = append(flags, layoutFlags...)
	flags = append(flags,
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

			if c.String("type") != "excel" {
				return &UserError{msg: "unsupported type: only 'excel' is supported now"}
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
				var oerr *testmtx.OpenFileError
				if errors.As(err, &oerr) {
					return &UserError{msg: oerr.Error()}
				}
				return fmt.Errorf("fail to create xlsx doc: %w", err)
			}

			sheetNames, err := getSheetNamesFunc(c, doc)
			if err != nil {
				return fmt.Errorf("fail to get sheet names: %w", err)
			}

			for _, sheetName := range sheetNames {
				docSheet, err := doc.GetSheet(sheetName)
				if err != nil {
					var nerr *testmtx.NotFoundError
					if errors.As(err, &nerr) {
						return &UserError{msg: nerr.Error()}
					}
					return fmt.Errorf("fail to get sheet: %w", err)
				}

				sheet, err := parser.Parse(docSheet)
				if err != nil {
					var perr *testmtx.ParseError
					if errors.As(err, &perr) {
						return &UserError{msg: fmt.Sprintf("parse error: %s", perr.Error())}
					}
					return fmt.Errorf("fail to parse sheet: %w", err)
				}

				if err := outputter.Output(c.String("out"), sheet); err != nil {
					var cerr *testmtx.CreateFileError
					if errors.As(err, &cerr) {
						return &UserError{msg: cerr.Error()}
					}
					var derr *testmtx.CreateDirError
					if errors.As(err, &derr) {
						return &UserError{msg: derr.Error()}
					}
					return fmt.Errorf("fail to output: %w", err)
				}
			}

			return nil
		},
	}, nil
}

type GetSheetNamesFlagAndFunc func() (
	flags []cli.Flag,
	getSheetNamesFunc func(*cli.Context, testmtx.Doc) ([]string, error),
)

type Formatter struct {
	Name    string
	NewFunc testmtx.NewFormatterFunc
}

type Layout struct {
	Name    string
	NewFunc testmtx.NewOutputterFunc
}

type getOutCommandOptions struct {
	getSheetNamesFlagAndFunc GetSheetNamesFlagAndFunc
	additionalSimpleValues   map[string]testmtx.ConvertValueFunc
	defaultPropLevel         int
	indentStr                string
	layouts                  []Layout
}

// When adding a new formatter, you also need to modify `testmtx.SimpleValue` interface.
// Adding a formatter is assumed to be added to `testmtx`, not as an option.
var formatters = []Formatter{
	{Name: "json", NewFunc: testmtx.NewJSONFormatter},
	{Name: "yaml", NewFunc: testmtx.NewYAMLFormatter},
}

var defaultNewOutCommandOptions = getOutCommandOptions{
	getSheetNamesFlagAndFunc: defaultGetSheetNamesFunc,
	additionalSimpleValues:   nil,
	defaultPropLevel:         10,
	indentStr:                "  ",
	layouts: []Layout{
		{Name: "1column-1case", NewFunc: testmtx.NewOneColumnOneCaseOutputter},
		{Name: "1sheet-1case", NewFunc: testmtx.NewOneSheetOneCaseOutputter},
	},
}

func defaultGetSheetNamesFunc() (
	flags []cli.Flag,
	getSheetNamesFunc func(*cli.Context, testmtx.Doc) ([]string, error),
) {
	return []cli.Flag{
			&cli.StringFlag{
				Name:    "sheet",
				Aliases: []string{"s"},
				Usage:   "input sheet name (required if --all-sheet is not specified)",
			},
			&cli.BoolFlag{
				Name:    "all-sheet",
				Aliases: []string{"all"},
				Usage:   "use all sheets in the file",
			},
		}, func(c *cli.Context, doc testmtx.Doc) ([]string, error) {
			if c.Bool("all-sheet") {
				sheetNames, err := doc.GetSheetNames()
				if err != nil {
					return nil, err
				}
				return sheetNames, nil
			}
			if sheetName := c.String("sheet"); sheetName != "" {
				return []string{sheetName}, nil
			}

			return nil, &UserError{msg: "sheet name is required (use --sheet or --all-sheet)"}
		}
}

type NewOutCommandOption func(*getOutCommandOptions)

func GetSheetNamesFlagAndFuncOption(fn GetSheetNamesFlagAndFunc) NewOutCommandOption {
	return func(o *getOutCommandOptions) {
		o.getSheetNamesFlagAndFunc = fn
	}
}

func AdditionalSimpleValues(convertValueFuncs map[string]testmtx.ConvertValueFunc) NewOutCommandOption {
	return func(o *getOutCommandOptions) {
		o.additionalSimpleValues = convertValueFuncs
	}
}

func DefaultPropLevel(level int) NewOutCommandOption {
	return func(o *getOutCommandOptions) {
		o.defaultPropLevel = level
	}
}

func IndentStr(indentStr string) NewOutCommandOption {
	return func(o *getOutCommandOptions) {
		o.indentStr = indentStr
	}
}

func Layouts(layouts []Layout) NewOutCommandOption {
	return func(o *getOutCommandOptions) {
		o.layouts = layouts
	}
}

func getParseFlagAndFunc(
	additionalSimpleValues map[string]testmtx.ConvertValueFunc,
	defaultPropLevel int,
) (
	flags []cli.Flag,
	newParserFunc func(c *cli.Context) (*testmtx.Parser, error),
	er error,
) {
	return []cli.Flag{
			&cli.IntFlag{
				Name:    "proplevel",
				Aliases: []string{"pl"},
				Value:   defaultPropLevel,
				Usage:   "property level (if you extend properties columns, then required)",
			},
		}, func(c *cli.Context) (*testmtx.Parser, error) {
			parser, err := testmtx.NewParser(
				testmtx.PropLevel(c.Int("proplevel")),
				testmtx.AdditionalSimpleValues(additionalSimpleValues),
			)
			if err != nil {
				return nil, fmt.Errorf("fail to create parser: %w", err)
			}

			return parser, nil
		}, nil
}

func getFormatFlagAndFunc(formatters []Formatter, indentStr string) (
	flags []cli.Flag,
	newFormatterFunc func(c *cli.Context) (testmtx.Formatter, error),
	er error,
) {
	switch len(formatters) {
	case 0:
		return nil, nil, errors.New("length of formatters must not be zero")
	case 1:
		return nil, func(c *cli.Context) (testmtx.Formatter, error) {
			f, err := formatters[0].NewFunc(testmtx.IndentStr(indentStr))
			if err != nil {
				return nil, fmt.Errorf("fail to create formatter: %w", err)
			}
			return f, nil
		}, nil
	}

	m := make(map[string]testmtx.NewFormatterFunc, len(formatters))
	usage := "output format type"
	for i, fs := range formatters {
		if _, ok := m[fs.Name]; ok {
			return nil, nil, fmt.Errorf("format name is duplicated (name: %q)", fs.Name)
		}
		m[fs.Name] = fs.NewFunc
		usage = appendUsageItem(usage, fs.Name, i, len(formatters))
	}

	return []cli.Flag{
			&cli.StringFlag{
				Name:  "format",
				Value: formatters[0].Name,
				Usage: usage,
			},
		}, func(c *cli.Context) (testmtx.Formatter, error) {
			fn, ok := m[c.String("format")]
			if !ok {
				return nil, errors.New("unsupportted format")
			}
			f, err := fn(testmtx.IndentStr(indentStr))
			if err != nil {
				return nil, fmt.Errorf("fail to create formatter: %w", err)
			}
			return f, nil
		}, nil
}

func getLayoutFlagAndFunc(layouts []Layout) (
	flags []cli.Flag,
	newOutputterFunc func(c *cli.Context, f testmtx.Formatter) (testmtx.Outputter, error),
	er error,
) {
	switch len(layouts) {
	case 0:
		return nil, nil, errors.New("length of layouts must not be zero")
	case 1:
		return nil, func(c *cli.Context, f testmtx.Formatter) (testmtx.Outputter, error) {
			return layouts[0].NewFunc(f), nil
		}, nil
	}

	m := make(map[string]testmtx.NewOutputterFunc, len(layouts))
	usage := "output file layout"
	for i, ls := range layouts {
		if _, ok := m[ls.Name]; ok {
			return nil, nil, fmt.Errorf("layout name is duplicated (name: %q)", ls.Name)
		}
		m[ls.Name] = ls.NewFunc
		usage = appendUsageItem(usage, ls.Name, i, len(layouts))
	}

	return []cli.Flag{
			&cli.StringFlag{
				Name:    "layout",
				Aliases: []string{"l"},
				Value:   layouts[0].Name,
				Usage:   usage,
			},
		}, func(c *cli.Context, f testmtx.Formatter) (testmtx.Outputter, error) {
			fn, ok := m[c.String("layout")]
			if !ok {
				return nil, errors.New("unsupportted layout")
			}
			return fn(f), nil
		}, nil
}
