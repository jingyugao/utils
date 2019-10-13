package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strings"

	"golang.org/x/tools/go/gcexportdata"
)

var newImporter = func(fset *token.FileSet) types.ImporterFrom {
	return gcexportdata.NewImporter(fset, make(map[string]*types.Package))
}

func main() {
	var dir string
	flag.StringVar(&dir, "d", "", "dir to scan")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of checkimpl:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	ifMap = map[string]*types.Interface{}
	implMap = map[string][]string{}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	fs := []*ast.File{}
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			fs = append(fs, f)
		}
	}
	config := &types.Config{
		Error:    func(error) {},
		Importer: newImporter(fset),
	}
	info := &types.Info{
		Types:  map[ast.Expr]types.TypeAndValue{},
		Defs:   map[*ast.Ident]types.Object{},
		Uses:   map[*ast.Ident]types.Object{},
		Scopes: map[ast.Node]*types.Scope{},
	}

	_, err = config.Check("", fset, fs, info)

	for k, v := range info.Defs {
		if v == nil {
			continue
		}
		t := info.TypeOf(k)
		if goType, ok := t.Underlying().(*types.Interface); ok {
			ifMap[v.Name()] = goType
		}
	}

	for k, v := range info.Defs {
		if v == nil {
			continue
		}
		t := info.TypeOf(k)
		checkImplements(v.Name(), t)
	}
	for tName, impls := range implMap {
		fmt.Printf("%s implement [%s]\n", tName, strings.Join(impls, ",\t"))
	}
}

var ifMap map[string]*types.Interface
var implMap map[string][]string

func checkImplements(name string, t types.Type) {
	if _, ok := t.Underlying().(*types.Struct); ok {
		checkImplementsAux("*"+name, types.NewPointer(t))
	}
	checkImplementsAux(name, t)
}

func checkImplementsAux(name string, t types.Type) {
	for k, itf := range ifMap {
		if k == name {
			continue
		}
		if types.Implements(t, itf) {
			implMap[name] = append(implMap[name], k)
			continue
		}
	}
}
