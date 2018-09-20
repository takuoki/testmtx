package sheet

import (
	"fmt"

	sheets "google.golang.org/api/sheets/v4"
)

// access to Google Sheets API

type client struct {
	sheetID string
	service *sheets.SpreadsheetsService
}

func new(authFile, spreadsheetID string) (*client, error) {
	srv, err := auth(authFile)
	if err != nil {
		return nil, err
	}

	return &client{
		sheetID: spreadsheetID,
		service: srv.Spreadsheets,
	}, nil
}

func (c *client) getSheetList() ([]string, error) {

	ss, err := c.service.Get(c.sheetID).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to get spreadsheet: %v", err)
	}

	sNames := []string{}
	for _, s := range ss.Sheets {
		sNames = append(sNames, s.Properties.Title)
	}

	return sNames, nil
}

func (c *client) getSheetData(sheetName string) ([][]interface{}, error) {

	resp, err := c.service.Values.Get(c.sheetID, sheetName).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve data from sheet: %v", err)
	}

	return resp.Values, nil
}
