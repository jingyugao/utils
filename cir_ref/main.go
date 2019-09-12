package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func main() {

	pkgMap, _ := parser.ParseDir(token.NewFileSet(), "./testdata", nil, parser.ParseComments)

	si := &sqlInject{}
	for k, v := range pkgMap {
		println(k)
		ast.Walk(si, v)
	}

	return
	allPath := CheckCircleRef(pkgMap)
	for _, path := range allPath {
		fmt.Printf("[%s] circle reference found\n", strings.Join(path, " -> "))
	}
}
