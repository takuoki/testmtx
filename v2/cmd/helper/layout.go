package helper

import (
	"errors"
	"fmt"

	"github.com/takuoki/testmtx/v2"
	cli "github.com/urfave/cli/v2"
)

type Layout struct {
	Name    string
	NewFunc testmtx.NewOutputterFunc
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
