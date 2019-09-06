package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

func printV(v interface{}) {
	r, e := json.Marshal(v)
	fmt.Println(string(r), e)
}

type V struct {
	pkgName       string
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
	doInDir(".", parseFile)

	for k := range struct2FieldKinds {
		s := []string{}
		if check(k, s) {
			println(k, strings.Join(s, ","))
		}
	}
}

func check(name string, stack []string) bool {
	return checkAux(name, stack)
}

func checkAux(cur string, stack []string) bool {
	for _, s := range stack {
		if s == cur {
			return true
		}
	}
	stack = append(stack, cur)

	vs, ok := struct2FieldKinds[cur]
	if !ok {
		return false
	}
	for _, v := range vs {
		if checkAux(v, stack) {
			return true
		}
	}

	return false
}

func hasChild(target string, cur string) bool {
	if vs, ok := struct2FieldKinds[cur]; ok {
		for _, v := range vs {
			if v == target {
				return true
			}
			if hasChild(target, v) {
				return true
			}
		}
	}
	return false
}

func parseFile(fileName string) {
	f, _ := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.ParseComments)
	v := V{}
	v.pkgName = f.Name.String()
	ast.Walk(&v, f)
	for k, v := range v.allStructInfo {
		struct2FieldKinds[k] = v
	}
}

func doInDir(dir string, fn func(string)) {
	flist, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, f := range flist {
		if f.IsDir() {
			doInDir(dir+"/"+f.Name(), fn)
		} else {
			if f.Name()[len(f.Name())-2:] == "go" {
				fn(dir + "/" + f.Name())
			}
		}
	}
}
