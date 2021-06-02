package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type config struct {
	ExcludedSheetNames []string `json:"excluded_sheet_names"`
	ExcludedSheetSet   map[string]struct{}
	SheetList          []sheetInfo `json:"sheet_list"`
	SheetAliasMap      map[string]string
}

type sheetInfo struct {
	Name    string `json:"name"`
	Alias   string `json:"alias"`
	SheetID string `json:"sheet_id"`
}

func newConfig() *config {
	return &config{
		ExcludedSheetNames: []string{},
		ExcludedSheetSet:   map[string]struct{}{},
		SheetList:          []sheetInfo{},
		SheetAliasMap:      map[string]string{},
	}
}

func (c *config) readConfig(configName string) error {

	jsonStr, err := ioutil.ReadFile(configName)
	if err != nil {
		return fmt.Errorf("config file not found (%s)", configName)
	}

	if err := json.Unmarshal(jsonStr, c); err != nil {
		return fmt.Errorf("invalid JSON format in config file (%s)", configName)
	}

	// ExcludedSheetSet
	c.ExcludedSheetSet = map[string]struct{}{}
	for _, sn := range c.ExcludedSheetNames {
		c.ExcludedSheetSet[sn] = struct{}{}
	}

	// SheetAliasMap
	c.SheetAliasMap = map[string]string{}
	for _, sa := range c.SheetList {
		c.SheetAliasMap[sa.Alias] = sa.SheetID
	}

	return nil
}

func (c *config) addExcludedSheet(s string) {
	if _, ok := c.ExcludedSheetSet[s]; !ok {
		c.ExcludedSheetNames = append(c.ExcludedSheetNames, s)
		c.ExcludedSheetSet[s] = struct{}{}
	}
}
