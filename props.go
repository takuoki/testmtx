package main

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
	"github.com/urfave/cli"

	"github.com/takuoki/testmtx/sheet"
)

var (
	lf              = []byte("\n")
	importMap       = map[string]string{}
	enumPropertyMap = map[string]*enumProperty{}
)

type enumProperty struct {
	name, value bool
}

func (e *enumProperty) valid() bool {
	return e.name && e.value
}

func init() {
	subCmdList = append(subCmdList, cli.Command{
		Name:  "prop",
		Usage: "output properties list based on Golang struct",
		Action: func(c *cli.Context) error {
			return action(c, &prop{})
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "file, f",
				Usage: "golang file witch target struct is written",
			},
			cli.StringFlag{
				Name:  "type, t",
				Usage: "target type name",
			},
		},
	})
}

type prop struct{}

func (p *prop) Run(c *cli.Context, _ *config) error {

	if c.String("file") == "" {
		return errors.New("no file name")
	}

	if c.String("type") == "" {
		return errors.New("no type name")
	}

	return p.Main(c.String("file"), c.String("type"))
}

func (p *prop) Main(file, tName string) error {

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

func (p *prop) outData(out io.Writer, d ast.Expr, i int) error {

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

func (p *prop) outObject(out io.Writer, t *ast.StructType, i int) error {
	out.Write([]byte(sheet.TypeObj))
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

func (p *prop) outArray(out io.Writer, t *ast.ArrayType, i int) error {
	out.Write([]byte(sheet.TypeAry))
	out.Write(lf)

	for j := 0; j < 2; j++ {
		p.outTab(out, i+1)
		out.Write([]byte("*"))
		p.outTab4Type(out, i+1)
		err := p.outData(out, t.Elt, i+1)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *prop) outIdent(out io.Writer, t *ast.Ident, i int) error {

	tName := ""
	switch t.Name {
	case "string", "rune":
		tName = sheet.TypeStr
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
		tName = sheet.TypeNum
	case "bool":
		tName = sheet.TypeBool
	default:
		if ep, ok := enumPropertyMap[t.Name]; ok && ep.valid() {
			tName = sheet.TypeStr
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

func (p *prop) outSelectorExpr(out io.Writer, t *ast.SelectorExpr, i int) error {

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
			tName = sheet.TypeStr
		}
	case "\"github.com/golang/protobuf/ptypes/wrappers\"":
		switch t.Sel.Name {
		case "StringValue", "BytesValue":
			tName = sheet.TypeStr
		case "Int32Value", "Int64Value", "UInt32Value", "UInt64Value", "FloatValue", "DoubleValue":
			tName = sheet.TypeNum
		case "BoolValue":
			tName = sheet.TypeBool
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

func (p *prop) outTab(out io.Writer, i int) {
	for j := 0; j < i; j++ {
		out.Write([]byte("\t"))
	}
}

func (p *prop) outTab4Type(out io.Writer, i int) {
	p.outTab(out, sheet.PropLevel-i)
}

func (p *prop) outKeyName(out io.Writer, tag *ast.BasicLit) error {

	if tag == nil {
		return errors.New("not found json tag")
	}

	pre := "json:\""
	s := strings.Index(tag.Value, pre) + len(pre)
	e := strings.Index(tag.Value[s:], "\"") + s
	out.Write([]byte(strings.Split(tag.Value[s:e], ",")[0]))

	return nil
}
