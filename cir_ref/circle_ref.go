package main

import (
	"fmt"
	"go/ast"
)

func CheckCircleRef(pkgMap map[string]*ast.Package) [][]string {
	struct2FieldKinds := map[string][]string{}
	v := structRef{}
	for pkgName, pkg := range pkgMap {
		v.pkgName = pkgName
		ast.Walk(&v, pkg)
	}

	for k, v := range v.allStructInfo {
		struct2FieldKinds[k] = v
	}
	return CheckCircle(struct2FieldKinds)
}

type structRef struct {
	pkgName       string // current walk pkg name
	allStructInfo map[string][]string
}

func (v *structRef) Visit(n ast.Node) ast.Visitor {
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
