package testmtx

// Sheet is ... TODO
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

func (v *vObject) firstProperty(cn casename, pn propname) bool {
	for _, n := range v.propertyNames {
		if n == pn {
			return true
		}
		if !v.properties[pn].isNil(cn) {
			return false
		}
	}
	return false
}

func (v *vObject) lastProperty(cn casename, pn propname) bool {
	check := false
	last := true
	for _, n := range v.propertyNames {
		if n == pn {
			check = true
			continue
		}
		if check {
			if !v.properties[n].isNil(cn) {
				last = false
				break
			}
		}
	}
	return last
}

type vArray struct {
	values   map[casename]bool
	elements []value
}

func (v *vArray) isNil(cn casename) bool {
	return !v.values[cn]
}

func (v *vArray) firstElement(cn casename, i int) bool {
	if i > len(v.elements) {
		return false
	}
	for j := 0; j < i; j++ {
		if !v.elements[j].isNil(cn) {
			return false
		}
	}
	return true
}

func (v *vArray) lastElement(cn casename, i int) bool {
	for i = i + 1; i < len(v.elements); i++ {
		if !v.elements[i].isNil(cn) {
			return false
		}
	}
	return true
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
