package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type Serializer struct {
	buf *bytes.Buffer
}

func (w *Serializer) Serializeint(v *int) error {
	return binary.Write(w.buf, binary.LittleEndian, v)
}

func main() {
	err := processFile("./foo.go")
	if err != nil {
		panic(err)
	}
}

func processFile(p string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, p, nil, 0)
	if err != nil {
		return err
	}

	for _, o := range f.Scope.Objects {
		// fmt.Printf("%#v\n", d)
		if o.Kind == ast.Typ {
			var p CodeParser
			err := p.processTypeSpec(o.Decl.(*ast.TypeSpec))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type CodeParser struct {
	state codeParserState

	fieldName string
	varName   string
}

type codeParserState int

const (
	stateIdle codeParserState = iota
	stateFieldArrayType
	stateFieldMapType
)

func (w *CodeParser) processTypeSpec(t *ast.TypeSpec) error {
	// fmt.Printf("%#v\n", t)
	structName := t.Name.Name

	fmt.Printf("func (w *Serializer) Serialize%s(v *%s) error {\n", structName, structName)

	err := w.processExpr(t.Type)
	if err != nil {
		return err
	}

	fmt.Println("}")

	return nil
}

func (w *CodeParser) processExpr(expr ast.Expr) error {
	// fmt.Printf("%#v\n", expr)

	switch e := expr.(type) {
	case *ast.StructType:
		return w.processStructType(e)
	case *ast.Ident:
		return w.processIdent(e)
	default:
		fmt.Printf("%#v\n", expr)
	}

	return nil
}

func (w *CodeParser) processIdent(e *ast.Ident) error {

	switch w.state {
	case stateFieldArrayType:
		eltTypeName := e.Name
		fmt.Printf("err = w.Serialize%s(&v.%s[i])\n", eltTypeName, w.fieldName)
		fmt.Println("if err != nil {")
		fmt.Println("return err")
		fmt.Println("}")
	case stateFieldMapType:
		typeName := e.Name
		fmt.Printf("err = w.Serialize%s(&%s)\n", typeName, w.varName)
		fmt.Println("if err != nil {")
		fmt.Println("return err")
		fmt.Println("}")
	}
	return nil
}

func (w *CodeParser) processStructType(st *ast.StructType) error {
	// fmt.Printf("%#v\n", st)

	for _, f := range st.Fields.List {
		err := w.processField(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *CodeParser) processField(f *ast.Field) error {
	// fmt.Printf("%#v\n", f)

	switch ft := f.Type.(type) {
	case *ast.Ident:
		return w.processFieldIdent(f, ft)
	case *ast.ArrayType:
		return w.processFieldArrayType(f, ft)
	case *ast.MapType:
		return w.processFieldMapType(f, ft)
	}

	// fmt.Printf("err = w.Serialize%s(v.%s)", fieldType, fieldName)

	return nil
}

func (w *CodeParser) processFieldIdent(f *ast.Field, ft *ast.Ident) error {
	fieldTypeName := ft.Name
	fieldName := f.Names[0].Name

	fmt.Printf("err = w.Serialize%s(&v.%s)\n", fieldTypeName, fieldName)
	fmt.Println("if err != nil {")
	fmt.Println("return err")
	fmt.Println("}")

	return nil
}

func (w *CodeParser) processFieldArrayType(f *ast.Field, ft *ast.ArrayType) error {
	w.state = stateFieldArrayType

	fieldName := f.Names[0].Name
	w.fieldName = fieldName

	fmt.Printf("for i := range v.%s {\n", fieldName)

	err := w.processExpr(ft.Elt)
	if err != nil {
		return err
	}

	fmt.Println("}")

	return nil
}

func (w *CodeParser) processFieldMapType(f *ast.Field, ft *ast.MapType) error {
	w.state = stateFieldMapType
	fieldName := f.Names[0].Name
	w.fieldName = fieldName

	fmt.Printf("for k, v := range v.%s {\n", fieldName)
	w.varName = "k"
	err := w.processExpr(ft.Key)
	if err != nil {
		return err
	}
	w.varName = "v"
	err = w.processExpr(ft.Value)
	if err != nil {
		return nil
	}
	fmt.Println("}")

	return nil
}
