package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type config struct {
	ExceptSheetNames []string `json:"except_sheet_names"`
	ExceptSheetSet   map[string]struct{}
	SheetList        []struct {
		Name    string `json:"name"`
		Alias   string `json:"alias"`
		SheetID string `json:"sheet_id"`
	} `json:"sheet_list"`
	SheetAliasMap map[string]string
}

func readConfig(configName string) (*config, error) {

	conf := &config{}

	jsonStr, err := ioutil.ReadFile(configName)
	if err != nil {
		return nil, fmt.Errorf("config file not found (%s)", configName)
	}

	if err := json.Unmarshal(jsonStr, conf); err != nil {
		return nil, fmt.Errorf("invalid JSON format in config file (%s)", configName)
	}

	// ExceptSheetSet
	conf.ExceptSheetSet = map[string]struct{}{}
	for _, sn := range conf.ExceptSheetNames {
		conf.ExceptSheetSet[sn] = struct{}{}
	}

	// SheetAliasMap
	conf.SheetAliasMap = map[string]string{}
	for _, sa := range conf.SheetList {
		conf.SheetAliasMap[sa.Alias] = sa.SheetID
	}

	return conf, nil
}
