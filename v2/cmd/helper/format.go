package helper

import (
	"errors"
	"fmt"

	"github.com/takuoki/testmtx/v2"
	cli "github.com/urfave/cli/v2"
)

type Formatter struct {
	Name    string
	NewFunc testmtx.NewFormatterFunc
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
