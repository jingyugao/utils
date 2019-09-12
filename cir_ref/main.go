package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type V struct {
	pkgName       string // current walk pkg name
	allStructInfo map[string][]string
}

func (v *V) Visit(n ast.Node) ast.Visitor {
	if v.allStructInfo == nil {
		v.allStructInfo = map[string][]string{}
	}
	if n == nil {
		return nil
	}

	if t, ok := n.(*ast.TypeSpec); ok {
		strName := t.Name.String()
		d := &DecodeStruct{
			pkgName: v.pkgName,
		}
		ast.Walk(d, t.Type)
		v.allStructInfo[v.pkgName+"."+strName] = d.load()
		return nil
	}
	return v
}

type DecodeStruct struct {
	pkgName   string
	KindNames map[string]bool
}

func (d *DecodeStruct) load() (fieldKinds []string) {
	for k := range d.KindNames {
		fieldKinds = append(fieldKinds, k)
	}
	return
}

func (d *DecodeStruct) addKindName(name string) {
	if d.KindNames == nil {
		d.KindNames = map[string]bool{}
	}
	d.KindNames[name] = true
}

func (d *DecodeStruct) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch v := n.(type) {
	case *ast.Field:
		d.Visit(v.Type)
		return nil
	case *ast.SelectorExpr:
		name := fmt.Sprintf("%s.%s", v.X, v.Sel)
		d.addKindName(name)
		return nil
	case *ast.Ident:
		d.addKindName(d.pkgName + "." + v.Name)
	case *ast.StarExpr:
		return d.Visit(v.X)
	case *ast.GenDecl:
		return nil
	}
	return d
}

var (
	struct2FieldKinds map[string][]string
)

func main() {
	struct2FieldKinds = map[string][]string{}
	ParseDir(".")
	allPath := CheckCircle(struct2FieldKinds)
	for _, path := range allPath {
		fmt.Printf("[%s] circle reference found\n", strings.Join(path, " -> "))
	}

}

func ParseDir(dir string) {
	pkgMap, _ := parser.ParseDir(token.NewFileSet(), dir, nil, parser.ParseComments)

	v := V{}
	for pkgName, pkg := range pkgMap {
		v.pkgName = pkgName
		ast.Walk(&v, pkg)
	}

	for k, v := range v.allStructInfo {
		struct2FieldKinds[k] = v
	}
}
