package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli"

	"github.com/takuoki/testmtx/sheet"
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
				Name:  "auth, a",
				Value: "credentials.json",
				Usage: "credential file for Google Sheets API",
			},
			cli.StringFlag{
				Name:  "sheet, s",
				Usage: "google spreadsheet id (mandatory)",
			},
			cli.StringFlag{
				Name:  "format, f",
				Value: "json",
				Usage: "output format (json, yaml)",
			},
			cli.StringFlag{
				Name:  "out, o",
				Value: "out",
				Usage: "output directory",
			},
			cli.IntFlag{
				Name:  "proplevel, pl",
				Value: 10,
				Usage: "properties level (if you extend properties columns, mandatory)",
			},
		},
	})
}

type output struct{}

type format interface {
	OutData(io.Writer, sheet.Data, sheet.Casename, int) error
	Extention() string
}

func (o *output) Run(c *cli.Context, conf *config) error {

	if c.String("sheet") == "" {
		return errors.New("no google spreadsheet id")
	}

	sheetID := c.String("sheet")
	if v, ok := conf.SheetAliasMap[sheetID]; ok {
		sheetID = v
	}
	ss, err := sheet.Get(c.String("auth"), sheetID, conf.ExceptSheetSet)
	if err != nil {
		return err
	}

	var f format
	switch c.String("format") {
	case "json":
		f = &jsonf{}
	case "yaml":
		f = &yamlf{}
	default:
		return fmt.Errorf("no such format (%s)", c.String("format"))
	}

	sheet.SetPropLevel(c.Int("proplevel"))

	if err := o.Main(ss, f, c.String("out")); err != nil {
		return err
	}

	fmt.Println("output completed successfully!")

	return nil
}

func (o *output) Main(ss []*sheet.Sheet, f format, outDir string) error {

	for _, s := range ss {
		for k, v := range s.DataMap {
			dir := fmt.Sprintf("%s/%s", outDir, k)
			if err := os.MkdirAll(dir, 0777); err != nil {
				return err
			}
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
