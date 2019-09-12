package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
)

func fnOpDB(fnName, pkgName string) bool {
	return fnName == "db" && pkgName == "Exec"
}

func log(i interface{}) {
	v, e := json.Marshal(i)
	fmt.Println(v, e)
}

type sqlInject struct {
}

func (v *sqlInject) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch v := n.(type) {
	case *ast.CallExpr:
		ast.Print(nil, v)
		return nil
	}
	return v
}

type funcWalk struct {
	pkg  string
	fn   string
	args string
}
