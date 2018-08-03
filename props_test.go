package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestOutKeyName(t *testing.T) {
	buf := &bytes.Buffer{}

	p := &prop{}
	p.outKeyName(buf, "`json:\"key3\"`")

	fmt.Println(buf)
	t.Error()
}
