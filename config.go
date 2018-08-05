package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type config struct {
	ExceptSheetNames []string `json:"except_sheet_names"`
	ExceptSheetSet   map[string]struct{}
	SheetAliasList   []struct {
		Alias   string `json:"alias"`
		SheetID string `json:"sheet_id"`
	} `json:"sheet_alias_list"`
	SheetAliasMap map[string]string
}

func readConfig(configName string) (*config, error) {

	conf := &config{}

	jsonStr, err := ioutil.ReadFile(configName)
	if err != nil {
		return nil, fmt.Errorf("not found config file (%s)", configName)
	}

	err = json.Unmarshal(jsonStr, conf)
	if err != nil {
		return nil, fmt.Errorf("something wrong in config file (%s)", configName)
	}

	// ExceptSheetSet
	conf.ExceptSheetSet = map[string]struct{}{}
	for _, sn := range conf.ExceptSheetNames {
		conf.ExceptSheetSet[sn] = struct{}{}
	}

	// SheetAliasMap
	conf.SheetAliasMap = map[string]string{}
	for _, sa := range conf.SheetAliasList {
		conf.SheetAliasMap[sa.Alias] = sa.SheetID
	}

	return conf, nil
}
