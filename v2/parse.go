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

func (p *Parser) propLevel(row DocRow) int {
	for i := p.propStartClm; i <= p.propEndClm; i++ {
		if row.Value(i) != "" {
			return i - p.propStartClm + 1
		}
	}
	// zero means no property title
	return 0
}

func (p *Parser) maxPropLevel() int {
	return p.propEndClm - p.propStartClm + 1
}

// Parse parses the sheet values to the sheet object.
// TODO: ParseErrorが返されることを明記する。
func (p *Parser) Parse(s DocSheet) (*Sheet, error) {

	sh := &Sheet{
		Name:        s.Name(),
		ColumnNames: []ColumnName{},
		Collections: map[PropName]Collection{},
	}

	rows := s.Rows()

	if s.Value(p.columnRow, p.columnStart) == "" {
		return nil, &ParseError{
			msg:       "first column name is empty",
			sheet:     pointer(s.Name()),
			rowNumber: pointer(rows[p.columnRow].Number()),
			clmLetter: pointer(clmconv.Itoa(p.columnStart)),
		}
	}

	// columns
	for ci := p.columnStart; ; ci++ {
		cn := s.Value(p.columnRow, ci)
		if cn == "" {
			break
		}
		for _, n := range sh.ColumnNames {
			if ColumnName(cn) == n {
				return nil, &ParseError{
					msg:       fmt.Sprintf("column name (%q) is duplicated", cn),
					sheet:     pointer(s.Name()),
					rowNumber: pointer(rows[p.columnRow].Number()),
					clmLetter: pointer(clmconv.Itoa(ci)),
				}
			}
		}

		sh.ColumnNames = append(sh.ColumnNames, ColumnName(strings.Replace(cn, " ", "_", -1)))
	}

	// properties
	for ri := p.dataStartRow; ri < len(rows); ri++ {
		if p.propLevel(rows[ri]) < 1 {
			continue
		}
		if lv := p.propLevel(rows[ri]); lv > 1 {
			return nil, &ParseError{
				msg:       "must not exist property that does not belong to the root property",
				sheet:     pointer(s.Name()),
				rowNumber: pointer(rows[ri].Number()),
				clmLetter: pointer(clmconv.Itoa(lv + p.propStartClm - 1)),
			}
		}
		pn := PropName(strings.Replace(rows[ri].Value(p.propStartClm), " ", "_", -1))
		if _, ok := sh.Collections[pn]; ok {
			return nil, &ParseError{
				msg:       fmt.Sprintf("root property name (%q) is duplicated", pn),
				sheet:     pointer(s.Name()),
				rowNumber: pointer(rows[ri].Number()),
				clmLetter: pointer(clmconv.Itoa(p.propStartClm)),
			}
		}
		var col Collection
		var err error
		col, ri, err = p.parseCollection(rows, ri, 1, sh.ColumnNames)
		if err != nil {
			return nil, &ParseError{
				msg:   fmt.Sprintf("fail to parse root property (%q)", pn),
				sheet: pointer(s.Name()),
				err:   err,
			}
		}
		sh.Collections[pn] = col
	}

	return sh, nil
}

func (p *Parser) parseCollection(rows []DocRow, ri, level int, cs []ColumnName) (Collection, int, error) {

	var col Collection
	var err error
	switch typ := rows[ri].Value(p.typeClm); typ {
	case typeObject:
		col, ri, err = p.parseObjectCollection(rows, ri, level, cs)
		if err != nil {
			return nil, 0, &ParseError{
				msg: "fail to create object collection",
				err: err,
			}
		}
	case typeArray:
		col, ri, err = p.parseArrayCollection(rows, ri, level, cs)
		if err != nil {
			return nil, 0, &ParseError{
				msg: "fail to create array collection",
				err: err,
			}
		}
	default:
		col, err = p.parseSimpleCollection(rows[ri], cs)
		if err != nil {
			return nil, 0, &ParseError{
				msg: fmt.Sprintf("fail to parse simple collection (type=%q)", typ),
				err: err,
			}
		}
	}

	return col, ri, nil
}

func (p *Parser) parseObjectCollection(rows []DocRow, ri, level int, cs []ColumnName) (*ObjectCollection, int, error) {

	implicitNils := map[ColumnName]bool{}
	explicitNils := map[ColumnName]bool{}

	for i := 0; i < len(cs); i++ {
		switch s := rows[ri].Value(p.columnStart + i); s {
		case "":
			implicitNils[cs[i]] = true
		case strNull:
			explicitNils[cs[i]] = true
		case strNew:
			// do nothing
		default:
			return nil, 0, &ParseError{
				msg:       fmt.Sprintf("invalid object value (%q)", s),
				rowNumber: pointer(rows[ri].Number()),
				clmLetter: pointer(clmconv.Itoa(p.columnStart + i)),
			}
		}
	}

	pNames := []PropName{}
	ps := map[PropName]Collection{}
	for ri = ri + 1; ri < len(rows); ri++ {

		if lv := p.propLevel(rows[ri]); lv < level+1 {
			ri--
			break
		} else if lv > level+1 {
			return nil, 0, &ParseError{
				msg:       "invalid level of object property",
				rowNumber: pointer(rows[ri].Number()),
				clmLetter: pointer(clmconv.Itoa(lv + p.propStartClm - 1)),
			}
		}

		pn := PropName(rows[ri].Value(p.propStartClm + level))
		pNames = append(pNames, pn)

		var err error
		ps[pn], ri, err = p.parseCollection(rows, ri, level+1, cs)
		if err != nil {
			return nil, 0, &ParseError{
				msg: "error occurred in object property",
				err: err,
			}
		}
	}

	return &ObjectCollection{
		ImplicitNils:  implicitNils,
		ExplicitNils:  explicitNils,
		PropertyNames: pNames,
		Properties:    ps,
	}, ri, nil
}

func (p *Parser) parseArrayCollection(rows []DocRow, ri, level int, cs []ColumnName) (*ArrayCollection, int, error) {

	implicitNils := map[ColumnName]bool{}
	explicitNils := map[ColumnName]bool{}

	for i := 0; i < len(cs); i++ {
		switch s := rows[ri].Value(p.columnStart + i); s {
		case "":
			implicitNils[cs[i]] = true
		case strNull:
			explicitNils[cs[i]] = true
		case strNew:
			// do nothing
		default:
			return nil, 0, &ParseError{
				msg:       fmt.Sprintf("invalid array value (%q)", s),
				rowNumber: pointer(rows[ri].Number()),
				clmLetter: pointer(clmconv.Itoa(p.columnStart + i)),
			}
		}
	}

	es := []Collection{}
	for ri = ri + 1; ri < len(rows); ri++ {

		if lv := p.propLevel(rows[ri]); lv < level+1 {
			ri--
			break
		} else if lv > level+1 {
			return nil, 0, &ParseError{
				msg:       "invalid level of array element",
				rowNumber: pointer(rows[ri].Number()),
				clmLetter: pointer(clmconv.Itoa(lv + p.propStartClm - 1)),
			}
		}

		var e Collection
		var err error
		e, ri, err = p.parseCollection(rows, ri, level+1, cs)
		if err != nil {
			return nil, 0, &ParseError{
				msg: "error occurred in array element",
				err: err,
			}
		}
		es = append(es, e)
	}

	return &ArrayCollection{
		ImplicitNils: implicitNils,
		ExplicitNils: explicitNils,
		Elements:     es,
	}, ri, nil
}

func (p *Parser) parseSimpleCollection(row DocRow, cs []ColumnName) (*SimpleCollection, error) {

	convertFunc, ok := p.convertSimpleValueFuncs[row.Value(p.typeClm)]
	if !ok {
		return nil, &ParseError{
			msg:       fmt.Sprintf("invalid type (%q)", row.Value(p.typeClm)),
			rowNumber: pointer(row.Number()),
			clmLetter: pointer(clmconv.Itoa(p.typeClm)),
		}
	}

	implicitNils := map[ColumnName]bool{}
	explicitNils := map[ColumnName]bool{}
	values := map[ColumnName]SimpleValue{}

	for i := 0; i < len(cs); i++ {
		switch s := row.Value(p.columnStart + i); s {
		case "":
			implicitNils[cs[i]] = true
		case strNull:
			explicitNils[cs[i]] = true
		default:
			v, err := convertFunc(s)
			if err != nil {
				return nil, &ParseError{
					msg:       "fail to convert simple value",
					rowNumber: pointer(row.Number()),
					clmLetter: pointer(clmconv.Itoa(p.columnStart + i)),
					err:       err,
				}
			}
			values[cs[i]] = v
		}
	}

	return &SimpleCollection{
		ImplicitNils: implicitNils,
		ExplicitNils: explicitNils,
		Values:       values,
	}, nil
}
