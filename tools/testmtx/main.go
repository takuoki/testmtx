package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

const version = "1.0.0"

var subCmdList = []cli.Command{}

type subCmd interface {
	Run(*cli.Context, *config) error
}

func main() {

	app := cli.NewApp()
	app.Name = "testmtx"
	app.Version = version
	app.Usage = "This tool is a test data generator using Google Spreadsheets."
	app.Commands = subCmdList
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "config file for testmtx",
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func action(c *cli.Context, sc subCmd) error {

	conf := &config{}
	if c.GlobalString("config") != "" {
		var err error
		conf, err = readConfig(c.GlobalString("config"))
		if err != nil {
			return err
		}
	}

	return sc.Run(c, conf)
}
