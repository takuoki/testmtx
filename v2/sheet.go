package testmtx

import "fmt"

type ColumnName string
type PropName string

// Sheet is a parsed sheet which matches testmtx format.
type Sheet struct {
	Name        string
	ColumnNames []ColumnName
	Collections map[PropName]Collection
}

type Collection interface {
	ImplicitNil(cn ColumnName) bool
	ExplicitNil(cn ColumnName) bool
}

type ObjectCollection struct {
	ImplicitNils  map[ColumnName]bool
	ExplicitNils  map[ColumnName]bool
	PropertyNames []PropName
	Properties    map[PropName]Collection
}

func newObjectCollection(cs []ColumnName, vs []string, pNames []PropName, ps map[PropName]Collection) (*ObjectCollection, error) {

	implicitNils := map[ColumnName]bool{}
	explicitNils := map[ColumnName]bool{}

	for i := 0; i < len(cs); i++ {
		switch str := vs[i]; str {
		case "":
			implicitNils[cs[i]] = true
		case strNull:
			explicitNils[cs[i]] = true
		case strNew:
			// do nothing
		default:
			return nil, fmt.Errorf("invalid object value (%q)", str)
		}
	}

	return &ObjectCollection{
		ImplicitNils:  implicitNils,
		ExplicitNils:  explicitNils,
		PropertyNames: pNames,
		Properties:    ps,
	}, nil
}

func (c *ObjectCollection) ImplicitNil(cn ColumnName) bool {
	return c.ImplicitNils[cn]
}

func (c *ObjectCollection) ExplicitNil(cn ColumnName) bool {
	return c.ExplicitNils[cn]
}

func (c *ObjectCollection) LastProperty(cn ColumnName, pn PropName) bool {
	for i := len(c.PropertyNames) - 1; i >= 0; i-- {
		if !c.Properties[c.PropertyNames[i]].ImplicitNil(cn) {
			return c.PropertyNames[i] == pn
		}
	}
	return false
}

type ArrayCollection struct {
	ImplicitNils map[ColumnName]bool
	ExplicitNils map[ColumnName]bool
	Elements     []Collection
}

func newArrayCollection(cs []ColumnName, vs []string, es []Collection) (*ArrayCollection, error) {

	implicitNils := map[ColumnName]bool{}
	explicitNils := map[ColumnName]bool{}

	for i := 0; i < len(cs); i++ {
		switch str := vs[i]; str {
		case "":
			implicitNils[cs[i]] = true
		case strNull:
			explicitNils[cs[i]] = true
		case strNew:
			// do nothing
		default:
			return nil, fmt.Errorf("invalid array value (%q)", str)
		}
	}

	return &ArrayCollection{
		ImplicitNils: implicitNils,
		ExplicitNils: explicitNils,
		Elements:     es,
	}, nil
}

func (c *ArrayCollection) ImplicitNil(cn ColumnName) bool {
	return c.ImplicitNils[cn]
}

func (c *ArrayCollection) ExplicitNil(cn ColumnName) bool {
	return c.ExplicitNils[cn]
}

func (c *ArrayCollection) LastElement(cn ColumnName, i int) bool {
	for j := len(c.Elements) - 1; j >= 0; j-- {
		if !c.Elements[j].ImplicitNil(cn) {
			return j == i
		}
	}
	return false
}

type SimpleCollection struct {
	ImplicitNils map[ColumnName]bool
	ExplicitNils map[ColumnName]bool
	Values       map[ColumnName]SimpleValue
}

func newSimpleCollection(cs []ColumnName, vs []string, fn ConvertValueFunc) (*SimpleCollection, error) {

	implicitNils := map[ColumnName]bool{}
	explicitNils := map[ColumnName]bool{}
	values := map[ColumnName]SimpleValue{}

	for i := 0; i < len(cs); i++ {
		switch vs[i] {
		case "":
			implicitNils[cs[i]] = true
		case strNull:
			explicitNils[cs[i]] = true
		default:
			v, err := fn(vs[i])
			if err != nil {
				return nil, err
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

func (c *SimpleCollection) ImplicitNil(cn ColumnName) bool {
	return c.ImplicitNils[cn]
}

func (c *SimpleCollection) ExplicitNil(cn ColumnName) bool {
	return c.ExplicitNils[cn]
}
