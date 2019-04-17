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
		Usage: "output properties list based on Golang struct",
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
				Usage: "properties level (if you extend properties columns, mandatory)",
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
		return errors.New("no file name")
	}

	if c.String("type") == "" {
		return errors.New("no type name")
	}

	pg, err := testmtx.NewPropGenerator(
		testmtx.PropLevel4Gen(c.Int("proplevel")),
		testmtx.RepeatCount(c.Int("repeated")),
	)
	if err != nil {
		return err
	}

	if err := pg.Generate(c.String("file"), c.String("type")); err != nil {
		return err
	}

	fmt.Println("\noutput completed successfully!")

	return nil
}
