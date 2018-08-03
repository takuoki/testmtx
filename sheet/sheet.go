package sheet

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Get is ...
func Get(spreadsheetID string) ([]*Sheet, error) {
	cli, err := new(spreadsheetID)
	if err != nil {
		return nil, err
	}

	sheetNames, err := cli.getSheetList()
	if err != nil {
		return nil, err
	}

	sheets := []*Sheet{}
	for _, sn := range sheetNames {
		d, err := cli.getSheetData(sn)
		if err != nil {
			return nil, err
		}
		s, err := parse(sn, d)
		if err != nil {
			return nil, err
		}
		sheets = append(sheets, s)
	}
	return sheets, nil
}

func parse(sheetName string, d [][]interface{}) (*Sheet, error) {

	if len(d) <= rDataStart {
		return nil, fmt.Errorf("sheet format is wrong (sheet=%s)", sheetName)
	}
	if len(d[rDataStart]) <= cCaseStart {
		return nil, fmt.Errorf("sheet format is wrong (sheet=%s)", sheetName)
	}

	s := &Sheet{
		Name:    sheetName,
		Cases:   []Casename{},
		DataMap: make(map[string]Data),
	}

	// cases
	for ci := cCaseStart; ci < len(d[rCase]); ci++ {
		cn := fmt.Sprintf("%s", d[rCase][ci])
		if cn == "" {
			break
		}
		for _, n := range s.Cases {
			if Casename(cn) == n {
				return nil, fmt.Errorf("case name is duplicated (sheet=%s, case=%s)", sheetName, cn)
			}
		}

		s.Cases = append(s.Cases, Casename(strings.Replace(cn, " ", "_", -1)))
	}

	// properties
	for ri := rDataStart; ri < len(d); ri++ {
		if level(d[ri]) == 1 {
			p := fmt.Sprintf("%s", d[ri][cPropStart])
			var pd Data
			var err error
			pd, ri, err = getData(d, ri, 1, s.Cases)
			if err != nil {
				return nil, fmt.Errorf("%s (sheet=%s)", err.Error(), sheetName)
			}
			s.DataMap[p] = pd
		}
	}

	return s, nil
}

func getData(d [][]interface{}, ri, l int, cases []Casename) (Data, int, error) {
	var dt Data
	switch tStr := d[ri][cType]; tStr {
	case TypeStr:
		v, err := getStringValues(d[ri], cases)
		if err != nil {
			return nil, ri, err
		}
		dt = &DString{
			Values: v,
		}
	case TypeNum:
		v, err := getNumValues(d[ri], cases)
		if err != nil {
			return nil, ri, err
		}
		dt = &DNum{
			Values: v,
		}
	case TypeBool:
		v, err := getBoolValues(d[ri], cases)
		if err != nil {
			return nil, ri, err
		}
		dt = &DBool{
			Values: v,
		}
	case TypeObj:
		v, err := getObjAryValues(d[ri], cases)
		if err != nil {
			return nil, ri, err
		}
		pns := []string{}
		pm := make(map[string]Data)
		for rj := ri + 1; rj < len(d); rj++ {
			if lv := level(d[rj]); lv < l+1 {
				break
			} else if lv > l+1 {
				return nil, rj, errors.New("the level of object child is wrong")
			}
			pn := fmt.Sprintf("%s", d[rj][cPropStart+l])
			pns = append(pns, pn)
			pm[pn], rj, err = getData(d, rj, l+1, cases)
			if err != nil {
				return nil, rj, err
			}
			ri = rj
		}
		dt = &DObject{
			Values:        v,
			PropertyNames: pns,
			Properties:    pm,
		}
	case TypeAry:
		v, err := getObjAryValues(d[ri], cases)
		if err != nil {
			return nil, ri, err
		}
		es := []Data{}
		for rj := ri + 1; rj < len(d); rj++ {
			if lv := level(d[rj]); lv < l+1 {
				break
			} else if lv > l+1 {
				return nil, rj, errors.New("the level of array child is wrong")
			}
			var tmpD Data
			tmpD, rj, err = getData(d, rj, l+1, cases)
			if err != nil {
				return nil, rj, err
			}
			es = append(es, tmpD)
			ri = rj
		}
		dt = &DArray{
			Values:   v,
			Elements: es,
		}
	default:
		return nil, ri, fmt.Errorf("no such type (%s)", tStr)
	}
	return dt, ri, nil
}

func level(rowData []interface{}) int {
	for i := cPropStart; i <= cPropEnd; i++ {
		if len(rowData) <= i {
			break
		}
		str := fmt.Sprintf("%s", rowData[i])
		if str != "" {
			return i - cPropStart + 1
		}
	}
	// zero means no property title
	return 0
}

func getObjAryValues(rowData []interface{}, cases []Casename) (map[Casename]bool, error) {

	m := make(map[Casename]bool)

	i := 0
	for ci := cCaseStart; ci < len(rowData); ci++ {
		if i >= len(cases) {
			break
		}
		str := fmt.Sprintf("%s", rowData[ci])
		switch str {
		case "":
			m[cases[i]] = false
		case strNew:
			m[cases[i]] = true
		default:
			return nil, fmt.Errorf("getObjAryValues: cannot convert array value (%s)", str)
		}

		i++
	}

	for ; i < len(cases); i++ {
		m[cases[i]] = false
	}

	return m, nil
}

func getStringValues(rowData []interface{}, cases []Casename) (map[Casename]*string, error) {

	m := make(map[Casename]*string)

	i := 0
	for ci := cCaseStart; ci < len(rowData); ci++ {
		if i >= len(cases) {
			break
		}
		str := fmt.Sprintf("%s", rowData[ci])
		switch str {
		case "":
			m[cases[i]] = nil
		case strEmpty:
			es := ""
			m[cases[i]] = &es
		default:
			m[cases[i]] = &str
		}

		i++
	}

	for ; i < len(cases); i++ {
		m[cases[i]] = nil
	}

	return m, nil
}

func getNumValues(rowData []interface{}, cases []Casename) (map[Casename]*string, error) {

	m := make(map[Casename]*string)

	i := 0
	for ci := cCaseStart; ci < len(rowData); ci++ {
		if i >= len(cases) {
			break
		}
		str := fmt.Sprintf("%s", rowData[ci])
		switch str {
		case "":
			m[cases[i]] = nil
		default:
			_, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return nil, fmt.Errorf("getNumValues: cannot convert int value (%s)", str)
			}
			m[cases[i]] = &str
		}

		i++
	}

	for ; i < len(cases); i++ {
		m[cases[i]] = nil
	}

	return m, nil
}

func getBoolValues(rowData []interface{}, cases []Casename) (map[Casename]*bool, error) {

	m := make(map[Casename]*bool)

	i := 0
	for ci := cCaseStart; ci < len(rowData); ci++ {
		if i >= len(cases) {
			break
		}
		str := fmt.Sprintf("%s", rowData[ci])
		switch str {
		case "":
			m[cases[i]] = nil
		case "true", "TRUE", "True":
			t := true
			m[cases[i]] = &t
		case "false", "FALSE", "False":
			f := false
			m[cases[i]] = &f
		default:
			return nil, fmt.Errorf("getBoolValues: cannot convert bool value (%s)", str)
		}

		i++
	}

	for ; i < len(cases); i++ {
		m[cases[i]] = nil
	}

	return m, nil
}
