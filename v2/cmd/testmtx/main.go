package main

import (
	"fmt"
	"log"
	"os"

	"github.com/takuoki/testmtx/v2/cmd/helper"
	cli "github.com/urfave/cli/v2"
)

const (
	appName = "testmtx"
	version = "2.0.0"
)

func main() {
	outCommand, err := helper.NewOutCommand()
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:    appName,
		Version: version,
		Usage:   "A test data generator using spreadsheet.",
		Commands: []*cli.Command{
			outCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(helper.ExtractErrorMsg(err))
		os.Exit(1)
	}
}
