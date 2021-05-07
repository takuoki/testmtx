package main

import (
	"errors"
	"fmt"

	"github.com/takuoki/testmtx"
	"github.com/urfave/cli"
)

func init() {
	subCmdList = append(subCmdList, cli.Command{
		Name:  "prop",
		Usage: "Outputs a property list based on Golang type",
		Action: func(c *cli.Context) error {
			return action(c, &prop{})
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "golang file in which the target type is defined (mandatory)",
			},
			cli.StringFlag{
				Name:  "type, t",
				Usage: "target type name (mandatory)",
			},
			cli.IntFlag{
				Name:  "proplevel, pl",
				Value: 10,
				Usage: "property level (if you extend properties columns, mandatory)",
			},
			cli.IntFlag{
				Name:  "repeated, r",
				Value: 2,
				Usage: "repeated count of array elements",
			},
		},
	})
}

type prop struct{}

func (p *prop) Run(c *cli.Context, _ *config) error {

	if c.String("file") == "" {
		return errors.New("please specify a file name")
	}

	if c.String("type") == "" {
		return errors.New("please specify a type name")
	}

	pg, err := testmtx.NewPropGenerator(
		testmtx.PropLevel4Gen(c.Int("proplevel")),
		testmtx.RepeatCount(c.Int("repeated")),
	)
	if err != nil {
		return fmt.Errorf("unable to create generator: %w", err)
	}

	if err := pg.Generate(c.String("file"), c.String("type")); err != nil {
		return fmt.Errorf("unable to generate a property list: %w", err)
	}

	fmt.Println("\ncomplete!")

	return nil
}
