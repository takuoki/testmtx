package helper

import (
	"fmt"

	"github.com/takuoki/testmtx/v2"
	cli "github.com/urfave/cli/v2"
)

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
