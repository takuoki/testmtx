package testmtx

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/takuoki/clmconv"
	"github.com/takuoki/gsheets/sheets"
)

// Property type in spreadsheet.
const (
	typeObj  = "object"
	typeAry  = "array"
	typeStr  = "string"
	typeNum  = "number"
	typeBool = "bool"
)

// String for special value.
const (
	strNothing = "*nothing"
	strNull    = "*null"
	strNew     = "*new"
	strEmpty   = "*empty"
)

// Parser is a struct to parse the sheet values to the sheet object.
// Create it using NewParser function.
type Parser struct {
	dataStartRow,
	propStartClm,
	propEndClm,
	typeClm,
	caseRow,
	caseStartClm int
}

// NewParser creates a new Parser.
// You can change some parameters of the Parser with ParseOption.
func NewParser(options ...ParseOption) (*Parser, error) {
	p := Parser{
		dataStartRow: 3,
		propStartClm: clmconv.MustAtoi("B"),
		propEndClm:   clmconv.MustAtoi("K"),
		typeClm:      clmconv.MustAtoi("L"),
		caseRow:      2,
		caseStartClm: clmconv.MustAtoi("M"),
	}
	for _, opt := range options {
		err := opt(&p)
		if err != nil {
			return nil, err
		}
	}
	return &p, nil
}

// ParseOption changes some parameters of the Parser.
type ParseOption func(*Parser) error

// PropLevel changes the property level on the spreadsheet.
func PropLevel(level int) ParseOption {
	return func(p *Parser) error {
		if level < 1 {
			return errors.New("property level should be positive value")
		}
		p.propEndClm = p.propStartClm + level - 1
		p.typeClm = p.propEndClm + 1
		p.caseStartClm = p.typeClm + 1
		return nil
	}
}

// Parse parses the sheet values to the sheet object.
func (p *Parser) Parse(s sheets.Sheet, sheetName string) (*Sheet, error) {

	if p == nil {
		return nil, errors.New("Parser is not initilized")
	}

	if s.Value(p.caseRow, p.caseStartClm) == "" {
		return nil, fmt.Errorf("invalid sheet format (sheet=%s)", sheetName)
	}

	sh := &Sheet{
		name:     sheetName,
		cases:    []casename{},
		valueMap: map[propname]value{},
	}

	// cases
	for ci := p.caseStartClm; ; ci++ {
		cn := s.Value(p.caseRow, ci)
		if cn == "" {
			break
		}
		for _, n := range sh.cases {
			if casename(cn) == n {
				return nil, fmt.Errorf("case name is duplicated (case=%s, sheet=%s)", cn, sheetName)
			}
		}

		sh.cases = append(sh.cases, casename(strings.Replace(cn, " ", "_", -1)))
	}

	// properties
	rows := s.Rows()
	for ri := p.dataStartRow; ri < len(rows); ri++ {
		if p.propLevel(rows[ri]) < 1 {
			continue
		}
		if p.propLevel(rows[ri]) > 1 {
			return nil, fmt.Errorf("must not exist property that does not belong to the root property (row=%d, sheet=%s)", ri, sheetName)
		}
		pn := propname(strings.Replace(rows[ri].Value(p.propStartClm), " ", "_", -1))
		if _, ok := sh.valueMap[pn]; ok {
			return nil, fmt.Errorf("root property name is duplicated (row=%d, sheet=%s)", ri, sheetName)
		}
		var val value
		var err error
		val, ri, err = p.getValues(rows, ri, 1, sh.cases)
		if err != nil {
			return nil, fmt.Errorf("%v, sheet=%s)", err, sheetName)
		}
		sh.valueMap[pn] = val
	}

	return sh, nil
}

func (p *Parser) maxPropLevel() int {
	return p.propEndClm - p.propStartClm + 1
}

func (p *Parser) propLevel(row sheets.Row) int {
	for i := p.propStartClm; i <= p.propEndClm; i++ {
		if row.Value(i) != "" {
			return i - p.propStartClm + 1
		}
	}
	// zero means no property title
	return 0
}

func (p *Parser) getValues(rows []sheets.Row, ri, l int, cases []casename) (value, int, error) {

	var val value
	switch rows[ri].Value(p.typeClm) {
	case typeObj:
		v, err := p.getObjAryValues(rows[ri], ri, cases)
		if err != nil {
			return nil, 0, err
		}
		var pns []propname
		pm := map[propname]value{}
		for ri = ri + 1; ri < len(rows); ri++ {
			if lv := p.propLevel(rows[ri]); lv < l+1 {
				ri--
				break
			} else if lv > l+1 {
				return nil, 0, fmt.Errorf("invalid level of object child (row=%d", ri)
			}
			pn := propname(rows[ri].Value(p.propStartClm + l))
			pns = append(pns, pn)
			pm[pn], ri, err = p.getValues(rows, ri, l+1, cases)
			if err != nil {
				return nil, 0, err
			}
		}
		val = &vObject{
			values:        v,
			propertyNames: pns,
			properties:    pm,
		}
	case typeAry:
		v, err := p.getObjAryValues(rows[ri], ri, cases)
		if err != nil {
			return nil, 0, err
		}
		var es []value
		for ri = ri + 1; ri < len(rows); ri++ {
			if lv := p.propLevel(rows[ri]); lv < l+1 {
				ri--
				break
			} else if lv > l+1 {
				return nil, 0, fmt.Errorf("invalid level of array child (row=%d", ri)
			}
			var tmpV value
			tmpV, ri, err = p.getValues(rows, ri, l+1, cases)
			if err != nil {
				return nil, 0, err
			}
			es = append(es, tmpV)
		}
		val = &vArray{
			values:   v,
			elements: es,
		}
	case typeStr:
		v, err := p.getStringValues(rows[ri], ri, cases)
		if err != nil {
			return nil, ri, err
		}
		val = &vString{
			values: v,
		}
	case typeNum:
		v, err := p.getNumValues(rows[ri], ri, cases)
		if err != nil {
			return nil, ri, err
		}
		val = &vNum{
			values: v,
		}
	case typeBool:
		v, err := p.getBoolValues(rows[ri], ri, cases)
		if err != nil {
			return nil, ri, err
		}
		val = &vBool{
			values: v,
		}
	default:
		return nil, ri, fmt.Errorf("invalid type (type=\"%s\", row=%d", rows[ri].Value(p.typeClm), ri)
	}
	return val, ri, nil
}

func (p *Parser) getObjAryValues(row sheets.Row, ri int, cases []casename) (map[casename]bool, error) {

	m := map[casename]bool{}

	for i := 0; i < len(cases); i++ {
		switch str := row.Value(p.caseStartClm + i); str {
		case strNull, "":
			m[cases[i]] = false
		case strNew:
			m[cases[i]] = true
		default:
			return nil, fmt.Errorf("unable to convert array value (value=\"%s\", row=%d", str, ri)
		}
	}

	return m, nil
}

func (p *Parser) getStringValues(row sheets.Row, ri int, cases []casename) (map[casename]*string, error) {

	m := map[casename]*string{}

	for i := 0; i < len(cases); i++ {
		switch str := row.Value(p.caseStartClm + i); str {
		case strNull, "":
			m[cases[i]] = nil
		case strEmpty:
			es := ""
			m[cases[i]] = &es
		default:
			m[cases[i]] = &str
		}
	}

	return m, nil
}

func (p *Parser) getNumValues(row sheets.Row, ri int, cases []casename) (map[casename]*string, error) {

	m := map[casename]*string{}

	for i := 0; i < len(cases); i++ {
		switch str := row.Value(p.caseStartClm + i); str {
		case strNull, "":
			m[cases[i]] = nil
		default:
			_, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return nil, fmt.Errorf("unable to convert int value (value=\"%s\", row=%d", str, ri)
			}
			m[cases[i]] = &str
		}
	}

	return m, nil
}

func (p *Parser) getBoolValues(row sheets.Row, ri int, cases []casename) (map[casename]*bool, error) {

	m := map[casename]*bool{}

	for i := 0; i < len(cases); i++ {
		switch str := row.Value(p.caseStartClm + i); str {
		case strNull, "":
			m[cases[i]] = nil
		case "true", "TRUE", "True":
			t := true
			m[cases[i]] = &t
		case "false", "FALSE", "False":
			f := false
			m[cases[i]] = &f
		default:
			return nil, fmt.Errorf("unable to convert bool value (value=\"%s\", row=%d", str, ri)
		}
	}

	return m, nil
}
