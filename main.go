package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

const (
	version = "0.1.0"
)

var subCmdList = []cli.Command{}

type subCmd interface {
	Run(*cli.Context) error
}

func main() {

	app := cli.NewApp()
	app.Name = "testmtx"
	app.Version = version
	app.Usage = "This tool is a test data generator using Google Spreadsheets."
	app.Commands = subCmdList

	// TODO: -auth
	// TODO: -config

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func action(c *cli.Context, sc subCmd) error {
	return sc.Run(c)
}
