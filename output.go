package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/takuoki/testmtx/sheet"
	"github.com/urfave/cli"
)

func init() {
	subCmdList = append(subCmdList, cli.Command{
		Name:  "out",
		Usage: "output test data files",
		Action: func(c *cli.Context) error {
			return action(c, &output{})
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "sheet, s",
				Usage: "google spreadsheet id",
			},
			cli.StringFlag{
				Name:  "format, f",
				Value: "json",
				Usage: "output format",
			},
		},
	})
}

type output struct{}

type format interface {
	OutData(io.Writer, sheet.Data, sheet.Casename, int) error
	Extention() string
}

func (o *output) Run(c *cli.Context) error {

	if c.String("sheet") == "" {
		return errors.New("no google spreadsheet id")
	}

	ss, err := sheet.Get(c.String("sheet"))
	if err != nil {
		return err
	}

	var f format
	switch c.String("format") {
	case "json":
		f = &jsonf{}
	default:
		return fmt.Errorf("no such format (%s)", c.String("format"))
	}

	return o.Main(ss, f)
}

func (o *output) Main(ss []*sheet.Sheet, f format) error {

	for _, s := range ss {
		for k, v := range s.DataMap {
			dir := fmt.Sprintf("out/%s", k)
			// TODO create dir
			for _, c := range s.Cases {
				err := o.outCase(f, s.Name, dir, c, v)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (o *output) outCase(f format, sheetName, dirName string, caseName sheet.Casename, d sheet.Data) error {
	file, err := os.Create(fmt.Sprintf(`./%s/%s_%s.%s`, dirName, sheetName, caseName, f.Extention()))
	if err != nil {
		return err
	}
	defer file.Close()

	return f.OutData(file, d, caseName, 0)
}
