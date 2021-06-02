package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
)

const version = "1.1.0"

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
		cli.StringFlag{
			Name:  "exclude, e",
			Usage: "sheet name to exclude (comma separated)",
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func action(c *cli.Context, sc subCmd) error {

	conf := newConfig()
	if c.GlobalString("config") != "" {
		if err := conf.readConfig(c.GlobalString("config")); err != nil {
			return err
		}
	}
	if c.GlobalString("exclude") != "" {
		ss := strings.Split(c.GlobalString("exclude"), ",")
		for _, s := range ss {
			conf.addExcludedSheet(s)
		}
	}

	return sc.Run(c, conf)
}
