package testmtx

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
)

var (
	lf              = []byte("\n")
	importMap       = map[string]string{}
	enumPropertyMap = map[string]*enumProperty{}

	repeated = 2 // overwritten by option
)

type enumProperty struct {
	name, value bool
}

func (e *enumProperty) valid() bool {
	return e.name && e.value
}

// PropGenerator is ... TODO
type PropGenerator struct {
	parser      *Parser
	repeatCount int
}

// NewPropGenerator creates a new PropGenerator.
// You can change some parameters of the PropGenerator with PropGenOption.
func NewPropGenerator(options ...PropGenOption) (*PropGenerator, error) {
	parser, err := NewParser()
	if err != nil {
		return nil, err
	}
	p := PropGenerator{
		parser:      parser,
		repeatCount: 3,
	}
	for _, opt := range options {
		err := opt(&p)
		if err != nil {
			return nil, err
		}
	}
	return &p, nil
}

// PropGenOption changes some parameters of the PropGenerator.
type PropGenOption func(*PropGenerator) error

// PropLevel4Gen changes the property level on the spreadsheet.
func PropLevel4Gen(level int) PropGenOption {
	return func(p *PropGenerator) error {
		if level < 1 {
			return errors.New("Property level should be positive value")
		}
		p.parser.propEndClm = p.parser.propStartClm + level - 1
		p.parser.typeClm = p.parser.propEndClm + 1
		p.parser.caseStartClm = p.parser.typeClm + 1
		return nil
	}
}

// RepeatCount changes the repeat count of the array properties.
func RepeatCount(c int) PropGenOption {
	return func(p *PropGenerator) error {
		if c < 1 {
			return errors.New("Repeat count should be positive value")
		}
		p.repeatCount = c
		return nil
	}
}

// Generate is ... TODO
func (p *PropGenerator) Generate(file, tName string) error {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		return err
	}

	// preparation
	for _, d := range f.Decls {
		if gd, ok := d.(*ast.GenDecl); ok {
			for _, s := range gd.Specs {
				if is, ok := s.(*ast.ImportSpec); ok {
					importMap[is.Name.Name] = is.Path.Value
				}
				// a type that satisfies the following conditions is an ENUM type generated by grpc
				// - the base type is "int32"
				// - variables with "_name", "_value" appended at the end are defined
				if ts, ok := s.(*ast.TypeSpec); ok {
					if id, ok := ts.Type.(*ast.Ident); ok {
						if id.Name == "int32" {
							enumPropertyMap[ts.Name.Name] = &enumProperty{}
						}
					}
				}
				if vs, ok := s.(*ast.ValueSpec); ok {
					if strings.Index(vs.Names[0].Name, "_name") > 0 {
						if ep, ok := enumPropertyMap[vs.Names[0].Name[0:len(vs.Names[0].Name)-5]]; ok {
							ep.name = true
						}
					}
					if strings.Index(vs.Names[0].Name, "_value") > 0 {
						if ep, ok := enumPropertyMap[vs.Names[0].Name[0:len(vs.Names[0].Name)-6]]; ok {
							ep.value = true
						}
					}
				}
			}
		}
	}

	// output
	for _, d := range f.Decls {
		if gd, ok := d.(*ast.GenDecl); ok {
			for _, s := range gd.Specs {
				if ts, ok := s.(*ast.TypeSpec); ok {
					if ts.Name.Name == tName {
						out := os.Stdout
						out.Write([]byte(strcase.ToSnake(tName)))
						p.outTab4Type(out, 0)
						return p.outData(out, ts.Type, 0)
					}
				}
			}
		}
	}

	return errors.New("no such type")
}

func (p *PropGenerator) outData(out io.Writer, d ast.Expr, i int) error {

	if i >= p.parser.maxLevel() {
		return fmt.Errorf("the type hierarchy exceeds the properties level. specify option '-proplevel' and re-execute")
	}

	var err error

	switch t := d.(type) {
	case *ast.StructType:
		err = p.outObject(out, t, i)
	case *ast.ArrayType:
		err = p.outArray(out, t, i)
	case *ast.Ident:
		err = p.outIdent(out, t, i)
	case *ast.SelectorExpr:
		err = p.outSelectorExpr(out, t, i)
	case *ast.StarExpr:
		err = p.outData(out, t.X, i)
	default:
		return fmt.Errorf("don't support type (%+v)", t)
	}

	return err
}

func (p *PropGenerator) outObject(out io.Writer, t *ast.StructType, i int) error {
	out.Write([]byte(typeObj))
	out.Write(lf)

	for _, f := range t.Fields.List {
		p.outTab(out, i+1)
		if err := p.outKeyName(out, f.Tag); err != nil {
			return err
		}
		p.outTab4Type(out, i+1)
		if err := p.outData(out, f.Type, i+1); err != nil {
			return err
		}
	}

	return nil
}

func (p *PropGenerator) outArray(out io.Writer, t *ast.ArrayType, i int) error {
	out.Write([]byte(typeAry))
	out.Write(lf)

	for j := 0; j < repeated; j++ {
		p.outTab(out, i+1)
		out.Write([]byte(fmt.Sprintf("* %d", j)))
		p.outTab4Type(out, i+1)
		err := p.outData(out, t.Elt, i+1)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PropGenerator) outIdent(out io.Writer, t *ast.Ident, i int) error {

	tName := ""
	switch t.Name {
	case "string", "rune":
		tName = typeStr
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
		tName = typeNum
	case "bool":
		tName = typeBool
	default:
		if ep, ok := enumPropertyMap[t.Name]; ok && ep.valid() {
			tName = typeStr
		} else {
			if t.Obj != nil {
				if ts, ok := t.Obj.Decl.(*ast.TypeSpec); ok {
					return p.outData(out, ts.Type, i)
				}
			}
			panic(fmt.Sprintf("don't expected type (%s)", t.Name))
		}
	}

	out.Write([]byte(tName))
	out.Write(lf)

	return nil
}

func (p *PropGenerator) outSelectorExpr(out io.Writer, t *ast.SelectorExpr, i int) error {

	var impPath string
	if x, ok := t.X.(*ast.Ident); ok {
		if pt, ok := importMap[x.Name]; ok {
			impPath = pt
		}
	}

	tName := ""
	switch impPath {
	case "\"github.com/golang/protobuf/ptypes/timestamp\"":
		if t.Sel.Name == "Timestamp" {
			tName = typeStr
		}
	case "\"github.com/golang/protobuf/ptypes/wrappers\"":
		switch t.Sel.Name {
		case "StringValue", "BytesValue":
			tName = typeStr
		case "Int32Value", "Int64Value", "UInt32Value", "UInt64Value", "FloatValue", "DoubleValue":
			tName = typeNum
		case "BoolValue":
			tName = typeBool
		}
	}

	if tName == "" {
		// if can't search, user modifies output data manually
		tName = fmt.Sprintf("<%s>", t.Sel.Name)
	}

	out.Write([]byte(tName))
	out.Write(lf)

	return nil
}

func (p *PropGenerator) outTab(out io.Writer, i int) {
	for j := 0; j < i; j++ {
		out.Write([]byte("\t"))
	}
}

func (p *PropGenerator) outTab4Type(out io.Writer, i int) {
	p.outTab(out, p.parser.maxLevel()-i)
}

func (p *PropGenerator) outKeyName(out io.Writer, tag *ast.BasicLit) error {

	if tag == nil {
		return errors.New("not found json tag")
	}

	pre := "json:\""
	s := strings.Index(tag.Value, pre) + len(pre)
	e := strings.Index(tag.Value[s:], "\"") + s
	out.Write([]byte(strings.Split(tag.Value[s:e], ",")[0]))

	return nil
}