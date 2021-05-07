package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/takuoki/gsheets"
	"github.com/takuoki/gsheets/sheets"
)

type gssDoc struct {
	ctx     context.Context
	client  *gsheets.Client
	sheetID string
}

func newGssDoc(sheetID, credentials string) (doc, error) {
	ctx := context.Background()
	ctx = gsheets.WithCache(ctx)
	client, err := gsheets.NewForCLI(ctx, credentials)
	if err != nil {
		return nil, fmt.Errorf("unable to create gsheets client: %w", err)
	}
	return &gssDoc{
		ctx:     ctx,
		client:  client,
		sheetID: sheetID,
	}, nil
}

func (d *gssDoc) GetSheetNames() ([]string, error) {
	if d == nil {
		return nil, errors.New("gssDoc is not initialized")
	}
	return d.client.GetSheetNames(d.ctx, d.sheetID)
}

func (d *gssDoc) GetSheet(sheetName string) (sheets.Sheet, error) {
	if d == nil {
		return nil, errors.New("gssDoc is not initialized")
	}
	return d.client.GetSheet(d.ctx, d.sheetID, sheetName)
}
