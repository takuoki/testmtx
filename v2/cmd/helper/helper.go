package helper

import (
	"errors"
	"fmt"

	"github.com/takuoki/testmtx/v2"
	cli "github.com/urfave/cli/v2"
)

func GetParseFlagAndFunc(convertValueFuncs map[string]testmtx.ConvertValueFunc) (
	flags []cli.Flag,
	newParserFunc func(c *cli.Context) (*testmtx.Parser, error),
	er error,
) {
	return []cli.Flag{
			&cli.IntFlag{
				Name:    "proplevel",
				Aliases: []string{"pl"},
				Value:   10,
				Usage:   "property level (if you extend properties columns, then required)",
			},
		}, func(c *cli.Context) (*testmtx.Parser, error) {
			parser, err := testmtx.NewParser(
				testmtx.PropLevel(c.Int("proplevel")),
				testmtx.AdditionalSimpleValues(convertValueFuncs),
			)
			if err != nil {
				return nil, fmt.Errorf("fail to create parser: %w", err)
			}

			return parser, nil
		}, nil
}

type Formatter struct {
	Name    string
	NewFunc testmtx.NewFormatterFunc
}

func GetFormatFlagAndFunc(formatters []Formatter, indentStr string) (
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
	usage := "output format type ("
	for i, fs := range formatters {
		if _, ok := m[fs.Name]; ok {
			return nil, nil, fmt.Errorf("format name is duplicated (name: %q)", fs.Name)
		}
		m[fs.Name] = fs.NewFunc
		if i < len(formatters)-2 {
			usage += fmt.Sprintf("%q, ", fs.Name)
		} else if i == len(formatters)-2 {
			usage += fmt.Sprintf("%q or ", fs.Name)
		} else {
			usage += fmt.Sprintf("%q", fs.Name)
		}
	}
	usage += ")"

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

type Outputter struct {
	Name    string
	NewFunc testmtx.NewOutputterFunc
}

func GetLayoutFlagAndFunc(outputters []Outputter) (
	flags []cli.Flag,
	newOutputterFunc func(c *cli.Context, f testmtx.Formatter) (testmtx.Outputter, error),
	er error,
) {
	switch len(outputters) {
	case 0:
		return nil, nil, errors.New("length of outputters must not be zero")
	case 1:
		return nil, func(c *cli.Context, f testmtx.Formatter) (testmtx.Outputter, error) {
			return outputters[0].NewFunc(f), nil
		}, nil
	}

	m := make(map[string]testmtx.NewOutputterFunc, len(outputters))
	usage := "output file layout ("
	for i, ls := range outputters {
		if _, ok := m[ls.Name]; ok {
			return nil, nil, fmt.Errorf("layout name is duplicated (name: %q)", ls.Name)
		}
		m[ls.Name] = ls.NewFunc
		if i < len(outputters)-2 {
			usage += fmt.Sprintf("%q, ", ls.Name)
		} else if i == len(outputters)-2 {
			usage += fmt.Sprintf("%q or ", ls.Name)
		} else {
			usage += fmt.Sprintf("%q", ls.Name)
		}
	}
	usage += ")"

	return []cli.Flag{
			&cli.StringFlag{
				Name:    "layout",
				Aliases: []string{"l"},
				Value:   outputters[0].Name,
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
