package sheet

import "errors"

const (
	// TypeObj is ...
	TypeObj = "object"

	// TypeAry is ...
	TypeAry = "array"

	// TypeStr is ...
	TypeStr = "string"

	// TypeNum is ...
	TypeNum = "number"

	// TypeBool is ...
	TypeBool = "bool"

	strNew   = "*new"
	strEmpty = "*empty"

	rCase      = 2
	cCaseStart = 12 // M

	rDataStart = 3
	cPropStart = 1  // B
	cType      = 11 // L
)

var (
	cPropEnd = 10 // K (default)
)

// SetPropLevel is ...
func SetPropLevel(level int) error {
	if level < 1 {
		return errors.New("SetPropLevel: level should be positive value")
	}

	cPropEnd = level + cPropStart - 1
	return nil
}

// GetPropLevel is ...
func GetPropLevel() int {
	return cPropEnd - cPropStart + 1
}

// Sheet is ...
type Sheet struct {
	Name    string
	Cases   []Casename
	DataMap map[string]Data
}

// Casename is ...
type Casename string

// Data is ...
type Data interface {
	IsNil(Casename) bool
}

// DObject is ...
type DObject struct {
	Values        map[Casename]bool
	PropertyNames []string
	Properties    map[string]Data
}

// IsNil is ...
func (d *DObject) IsNil(c Casename) bool {
	return !d.Values[c]
}

// FirstProperty is ...
func (d *DObject) FirstProperty(c Casename, propName string) bool {
	for _, pn := range d.PropertyNames {
		if pn == propName {
			return true
		}
		if !d.Properties[pn].IsNil(c) {
			return false
		}
	}
	return false
}

// LastProperty is ...
func (d *DObject) LastProperty(c Casename, propName string) bool {
	check := false
	last := true
	for _, pn := range d.PropertyNames {
		if pn == propName {
			check = true
			continue
		}
		if check {
			if !d.Properties[pn].IsNil(c) {
				last = false
				break
			}
		}
	}
	return last
}

// DArray is ...
type DArray struct {
	Values   map[Casename]bool
	Elements []Data
}

// IsNil is ...
func (d *DArray) IsNil(c Casename) bool {
	return !d.Values[c]
}

// FirstElement is ...
func (d *DArray) FirstElement(c Casename, i int) bool {
	if i > len(d.Elements) {
		return false
	}
	for j := 0; j < i; j++ {
		if !d.Elements[j].IsNil(c) {
			return false
		}
	}
	return true
}

// LastElement is ...
func (d *DArray) LastElement(c Casename, i int) bool {
	for i = i + 1; i < len(d.Elements); i++ {
		if !d.Elements[i].IsNil(c) {
			return false
		}
	}
	return true
}

// DString is ...
type DString struct {
	Values map[Casename]*string
}

// IsNil is ...
func (d *DString) IsNil(c Casename) bool {
	return d.Values[c] == nil
}

// DNum is ...
type DNum struct {
	Values map[Casename]*string
}

// IsNil is ...
func (d *DNum) IsNil(c Casename) bool {
	return d.Values[c] == nil
}

// DBool is ...
type DBool struct {
	Values map[Casename]*bool
}

// IsNil is ...
func (d *DBool) IsNil(c Casename) bool {
	return d.Values[c] == nil
}
