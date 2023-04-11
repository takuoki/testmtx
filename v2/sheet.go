package testmtx

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

func (c *SimpleCollection) ImplicitNil(cn ColumnName) bool {
	return c.ImplicitNils[cn]
}

func (c *SimpleCollection) ExplicitNil(cn ColumnName) bool {
	return c.ExplicitNils[cn]
}
