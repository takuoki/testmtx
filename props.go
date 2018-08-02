package main

import (
	"fmt"

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
				Usage: "golang file witch target struct is written",
			},
			cli.StringFlag{
				Name:  "struct, s",
				Usage: "target struct name",
			},
		},
	})
}

type prop struct{}

func (p *prop) Run(c *cli.Context) error {
	fmt.Println("not implemented")
	return nil
}
