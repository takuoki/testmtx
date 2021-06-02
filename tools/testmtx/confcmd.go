package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func init() {
	subCmdList = append(subCmdList, cli.Command{
		Name:  "conf",
		Usage: "Outputs configuration file",
		Action: func(c *cli.Context) error {
			return action(c, &conf{w: os.Stdout})
		},
	})
}

type conf struct {
	w io.Writer
}

func (c *conf) Run(ctx *cli.Context, conf *config) error {

	if ctx.GlobalString("config") == "" {
		return errors.New("please specify a configuration file")
	}

	fmt.Fprintln(c.w, "# Excluded Sheet Names")
	for _, s := range conf.ExcludedSheetNames {
		fmt.Fprintf(c.w, "- %s\n", s)
	}

	fmt.Fprintln(c.w, "\n# Sheet List")
	table := tablewriter.NewWriter(c.w)
	table.SetHeader([]string{"Name", "Alias", "Spreadsheet ID"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("-")
	table.SetBorder(false)
	for _, s := range conf.SheetList {
		table.Append([]string{s.Name, s.Alias, s.SheetID})
	}
	table.Render()

	return nil
}
