package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func main() {
	fset := token.NewFileSet()
	si := &sqlInject{fset: fset}

	file, err := parser.ParseFile(fset, "./testdata/const.go", nil, 0)
	// ast.Walk(si, file)

	ast.Print(fset, file)
	return

	pkgMap, err := parser.ParseDir(fset, "./testdata", nil, 0)
	// pkgMap, err := parser.ParseDir(fset, "/Users/gao/Code/stat/model", nil, 0)
	if err != nil {
		panic(err)
	}
	for _, v := range pkgMap {
		f := "/Users/gao/Code/stat/model/mstdealercount.go"
		ast.Walk(si, v.Files[f])
		// ast.Walk(si, v)
	}

	_ = si
	return
	allPath := CheckCircleRef(pkgMap)
	for _, path := range allPath {
		fmt.Printf("[%s] circle reference found\n", strings.Join(path, " -> "))
	}
}
