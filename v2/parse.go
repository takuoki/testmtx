package testmtx

import (
	"errors"
	"fmt"
	"strings"

	"github.com/takuoki/clmconv"
)

// Parser is a struct to parse the sheet values to the sheet object.
// Create it using NewParser function.
type Parser struct {
	dataStartRow,
	propStartClm,
	propEndClm,
	typeClm,
	columnRow,
	columnStart int
	convertSimpleValueFuncs map[string]ConvertValueFunc
}

// NewParser creates a new Parser.
// You can change some parameters of the Parser with ParseOption.
func NewParser(options ...ParseOption) (*Parser, error) {
	p := Parser{
		dataStartRow:            3,
		propStartClm:            clmconv.MustAtoi("B"),
		propEndClm:              clmconv.MustAtoi("K"),
		typeClm:                 clmconv.MustAtoi("L"),
		columnRow:               2,
		columnStart:             clmconv.MustAtoi("M"),
		convertSimpleValueFuncs: defaultConvertSimpleValueFuncs,
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

// PropLevel changes the property level on the sheet.
func PropLevel(level int) ParseOption {
	return func(p *Parser) error {
		if level < 1 {
			return errors.New("property level should be positive value")
		}
		p.propEndClm = p.propStartClm + level - 1
		p.typeClm = p.propEndClm + 1
		p.columnStart = p.typeClm + 1
		return nil
	}
}

// AdditionalSimpleValues adds simple values to default simple value list.
func AdditionalSimpleValues(convertValueFuncs map[string]ConvertValueFunc) ParseOption {
	return func(p *Parser) error {
		for k, v := range convertValueFuncs {
			if _, ok := p.convertSimpleValueFuncs[k]; ok {
				return fmt.Errorf("type name %q is duplicated", k)
			}
			p.convertSimpleValueFuncs[k] = v
		}
		return nil
	}
}

// Parse parses the sheet values to the sheet object.
func (p *Parser) Parse(s DocSheet) (*Sheet, error) {

	if p == nil {
		return nil, errors.New("parser is not initilized")
	}

	if s.Value(p.columnRow, p.columnStart) == "" {
		return nil, fmt.Errorf("invalid sheet format (sheet=%s)", s.Name())
	}

	sh := &Sheet{
		Name:        s.Name(),
		ColumnNames: []ColumnName{},
		Collections: map[PropName]Collection{},
	}

	// columns
	for ci := p.columnStart; ; ci++ {
		cn := s.Value(p.columnRow, ci)
		if cn == "" {
			break
		}
		for _, n := range sh.ColumnNames {
			if ColumnName(cn) == n {
				return nil, fmt.Errorf("column name is duplicated (name=%s, sheet=%s)", cn, s.Name())
			}
		}

		sh.ColumnNames = append(sh.ColumnNames, ColumnName(strings.Replace(cn, " ", "_", -1)))
	}

	// properties
	rows := s.Rows()
	for ri := p.dataStartRow; ri < len(rows); ri++ {
		if p.propLevel(rows[ri]) < 1 {
			continue
		}
		if p.propLevel(rows[ri]) > 1 {
			return nil, fmt.Errorf("must not exist property that does not belong to the root property (row=%d, sheet=%s)", ri, s.Name())
		}
		pn := PropName(strings.Replace(rows[ri].Value(p.propStartClm), " ", "_", -1))
		if _, ok := sh.Collections[pn]; ok {
			return nil, fmt.Errorf("root property name is duplicated (row=%d, sheet=%s)", ri, s.Name())
		}
		var col Collection
		var err error
		col, ri, err = p.getValues(rows, ri, 1, sh.ColumnNames)
		if err != nil {
			// TODO: エラーの返し方
			return nil, fmt.Errorf("%v, sheet=%s)", err, s.Name())
		}
		sh.Collections[pn] = col
	}

	return sh, nil
}

func (p *Parser) getValues(rows []DocRow, ri, l int, cs []ColumnName) (Collection, int, error) {

	var col Collection
	var err error
	switch typ := rows[ri].Value(p.typeClm); typ {
	case typeObject:
		vs := p.getValueStrings(rows[ri], cs)

		pNames := []PropName{}
		ps := map[PropName]Collection{}
		for ri = ri + 1; ri < len(rows); ri++ {
			if lv := p.propLevel(rows[ri]); lv < l+1 {
				ri--
				break
			} else if lv > l+1 {
				// TODO: エラーの返し方（行番号でのラップが必要）
				return nil, 0, fmt.Errorf("invalid level of object child (row=%d", ri)
			}
			pn := PropName(rows[ri].Value(p.propStartClm + l))
			pNames = append(pNames, pn)
			ps[pn], ri, err = p.getValues(rows, ri, l+1, cs)
			if err != nil {
				// TODO: エラーの返し方（行番号でのラップが必要）
				return nil, 0, err
			}
		}

		col, err = newObjectCollection(cs, vs, pNames, ps)
		if err != nil {
			return nil, 0, fmt.Errorf("error occurred (row=%d): %w", ri, err)
		}
	case typeArray:
		vs := p.getValueStrings(rows[ri], cs)

		es := []Collection{}
		for ri = ri + 1; ri < len(rows); ri++ {
			if lv := p.propLevel(rows[ri]); lv < l+1 {
				ri--
				break
			} else if lv > l+1 {
				// TODO: エラーの返し方（行番号でのラップが必要）
				return nil, 0, fmt.Errorf("invalid level of array child (row=%d", ri)
			}
			var tmpV Collection
			tmpV, ri, err = p.getValues(rows, ri, l+1, cs)
			if err != nil {
				// TODO: エラーの返し方（行番号でのラップが必要）
				return nil, 0, err
			}
			es = append(es, tmpV)
		}

		col, err = newArrayCollection(cs, vs, es)
		if err != nil {
			return nil, 0, fmt.Errorf("error occurred (row=%d): %w", ri, err)
		}
	default:
		convertFunc, ok := p.convertSimpleValueFuncs[typ]
		if !ok {
			// TODO: エラーの返し方（行番号でのラップが必要）
			return nil, ri, fmt.Errorf("invalid type (type=\"%s\", row=%d", rows[ri].Value(p.typeClm), ri)
		}

		col, err = newSimpleCollection(cs, p.getValueStrings(rows[ri], cs), convertFunc)
		if err != nil {
			return nil, 0, fmt.Errorf("error occurred (row=%d): %w", ri, err)
		}
	}

	return col, ri, nil
}

func (p *Parser) getValueStrings(row DocRow, cs []ColumnName) []string {
	values := []string{}
	for i := 0; i < len(cs); i++ {
		values = append(values, row.Value(p.columnStart+i))
	}
	return values
}

func (p *Parser) maxPropLevel() int {
	return p.propEndClm - p.propStartClm + 1
}

func (p *Parser) propLevel(row DocRow) int {
	for i := p.propStartClm; i <= p.propEndClm; i++ {
		if row.Value(i) != "" {
			return i - p.propStartClm + 1
		}
	}
	// zero means no property title
	return 0
}
