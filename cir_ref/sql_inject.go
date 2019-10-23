package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/token"
)

func log(i interface{}) {
	v, e := json.Marshal(i)
	fmt.Println(string(v), e)
}

type sqlInject struct {
	fset      *token.FileSet
	inDBQuery bool
}

func (si *sqlInject) Visit(n ast.Node) ast.Visitor {
	defer func() {
		// if err := recover(); err != nil {
		// 	p := n.Pos()
		// 	fmt.Printf("[%s] \n", si.fset.Position(p))
		// 	os.Exit(1)
		// }
	}()
	if n == nil {
		return nil
	}

	switch v := n.(type) {
	case *ast.CallExpr:
		inDBQuery := new(checkFnDB)
		inDBQuery.formatArgIdx = -1
		ast.Walk(inDBQuery, v.Fun)
		if inDBQuery.formatArgIdx == -1 {
			return nil
		}
		idx := inDBQuery.formatArgIdx
		if len(v.Args) <= idx {
			println("idx out of range")
			return nil
		}

		if !checkVarIsConst(v.Args[idx]) {
			p := n.Pos()
			fmt.Printf("%s \n", si.fset.Position(p))
		}

		return nil
	}
	return si
}

type checkFnDB struct {
	formatArgIdx int
}

func (c *checkFnDB) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch v := n.(type) {
	case *ast.Ident:
		switch v.String() {
		case "Select", "Get":
			c.formatArgIdx = 1
		case "Query", "Queryx", "Exec":
			c.formatArgIdx = 0
		case "SelectContext", "GetContext":
			c.formatArgIdx = 2
		case "QueryContext", "QueryxContext", "ExecContext":
			c.formatArgIdx = 1
		}

		return nil
	}
	return c
}

func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

func isPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}

func checkVarIsConst(n ast.Node) bool {
	switch v := n.(type) {
	case *ast.Ident:
		if v.Obj != nil && v.Obj.Kind == ast.Con {
			return true
		}
		if v.Obj == nil {
			return true
		}
		if v.Obj.Decl == nil {
			return true
		}
		if assign, ok := v.Obj.Decl.(*ast.AssignStmt); ok {
			return checkVarIsConst(assign.Rhs[0])
		}
		ast.Print(nil, n)
	case *ast.BasicLit:
		return true
	case *ast.BinaryExpr:
		return checkVarIsConst(v.X) && checkVarIsConst(v.Y)
	case *ast.CallExpr:
		if sel, ok := v.Fun.(*ast.SelectorExpr); ok && isIdent(sel.X, "builder") {
			return true
		}
		if sel, ok := v.Fun.(*ast.SelectorExpr); ok && isIdent(sel.X, "reflect") {
			return true
		}
		if sel, ok := v.Fun.(*ast.SelectorExpr); ok && isIdent(sel.Sel, "NamedSelectSQL") {
			return true
		}
		if isIdent(v.Fun, "namedSelectSQL") {
			return true
		}
		for _, arg := range v.Args {
			if !checkVarIsConst(arg) {
				return false
			}
		}
		return true
	}

	return false
}
