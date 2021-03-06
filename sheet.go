// Package testmtx helps you to create test data files with Google Spreadsheets.
// Once you create test cases as matrix on Google Spreadsheets,
// this tool generates test data like JSON based on the data you input.
// Using testmtx, you can get advantages of `completeness`, `readability` and `consistency` for testing.
//
// This package is just liblary.
// If what you want is just to use, see the standard tool below which uses this package.
// - github.com/takuoki/testmtx/tools/testmtx
package testmtx

// Sheet is a parsed sheet which matches testmtx format.
type Sheet struct {
	name     string
	cases    []casename
	valueMap map[propname]value
}

type casename string
type propname string

type value interface {
	isNil(cn casename) bool
}

type vObject struct {
	values        map[casename]bool
	propertyNames []propname
	properties    map[propname]value
}

func (v *vObject) isNil(cc casename) bool {
	return !v.values[cc]
}

func (v *vObject) lastProperty(cn casename, pn propname) bool {
	for i := len(v.propertyNames) - 1; i >= 0; i-- {
		if !v.properties[v.propertyNames[i]].isNil(cn) {
			return v.propertyNames[i] == pn
		}
	}
	return false
}

type vArray struct {
	values   map[casename]bool
	elements []value
}

func (v *vArray) isNil(cn casename) bool {
	return !v.values[cn]
}

func (v *vArray) lastElement(cn casename, i int) bool {
	for j := len(v.elements) - 1; j >= 0; j-- {
		if !v.elements[j].isNil(cn) {
			return j == i
		}
	}
	return false
}

type vString struct {
	values map[casename]*string
}

func (v *vString) isNil(cn casename) bool {
	return v.values[cn] == nil
}

type vNum struct {
	values map[casename]*string
}

func (v *vNum) isNil(cn casename) bool {
	return v.values[cn] == nil
}

type vBool struct {
	values map[casename]*bool
}

func (v *vBool) isNil(cn casename) bool {
	return v.values[cn] == nil
}
