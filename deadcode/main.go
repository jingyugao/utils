package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io/ioutil"
	"os"
	"sort"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

func main() {
	cfg := &packages.Config{}
	cfg.Fset = token.NewFileSet()
	cfg.Mode = packages.LoadAllSyntax

	pkgs, err := packages.Load(cfg, os.Args[1])
	if err != nil {
		panic(err)
	}
	prog := &prog{
		pkgs: pkgs,
		decl: map[string]ast.Node{},
		used: map[string]bool{},
		fs:   cfg.Fset,
	}
	prog.run()
	return

}

type prog struct {
	pkgs []*packages.Package
	decl map[string]ast.Node
	used map[string]bool
	fs   *token.FileSet
}

type Package struct {
	p    *packages.Package
	prog *prog
}

func (prog *prog) run() {
	for _, pkg := range prog.pkgs {
		prog.doPackage(pkg)
	}
	// reports.
	reports := Reports(nil)
	for name, node := range prog.decl {
		if !prog.used[name] {
			reports = append(reports, Report{node.Pos(), name})
		}
	}

	sort.Sort(reports)
	for _, report := range reports {
		fmt.Printf("%s: %s is unused\n", prog.fs.Position(report.pos), report.name)
	}

	for _, pkg := range prog.pkgs {
		prog.doTrim(pkg)
	}
}

func (p *prog) pre(c *astutil.Cursor) bool {
	n := c.Node()
	if d, ok := n.(*ast.FuncDecl); ok {
		if !p.used[d.Name.Name] {
			c.Delete()
			return false
		}
	}
	return true
}
func (prog *prog) doTrim(pkg *packages.Package) {
	for _, f := range pkg.Syntax {
		astutil.Apply(f, prog.pre, nil)
		var buf bytes.Buffer
		err := format.Node(&buf, prog.fs, f)
		if err != nil {
			panic(err)
		}
		fName := prog.fs.Position(f.Pos()).Filename
		err = ioutil.WriteFile(fName, buf.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
	}
}
func (prog *prog) doPackage(pkg *packages.Package) {

	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			switch n := decl.(type) {
			case *ast.GenDecl:
				// var, const, types
				// for _, spec := range n.Specs {
				// 	switch s := spec.(type) {
				// 	case *ast.ValueSpec:
				// 		// constants and variables.
				// 		for _, name := range s.Names {
				// 			p.prog.decl[name.Name] = n
				// 		}
				// 	case *ast.TypeSpec:
				// 		// type definitions.
				// 		p.prog.decl[s.Name.Name] = n
				// 	}
				// }
			case *ast.FuncDecl:
				// function declarations
				// TODO(remy): do methods
				if n.Recv == nil {
					prog.decl[n.Name.Name] = n
				}
			}
		}
	}
	// init() and _ are always used
	prog.used["init"] = true
	prog.used["_"] = true
	if pkg.Name != "main" {
		// exported names are marked used for non-main packages.
		// for name := range p.prog.decl {
		// 	if ast.IsExported(name) {
		// 		p.prog.used[p.p.Name+"."+name] = true
		// 	}
		// }
	} else {
		// in main programs, main() is called.
		prog.used["main"] = true
	}
	for _, file := range pkg.Syntax {
		// walk file looking for used nodes.
		ast.Walk(&Package{pkg, prog}, file)
	}
}

type Report struct {
	pos  token.Pos
	name string
}
type Reports []Report

func (l Reports) Len() int           { return len(l) }
func (l Reports) Less(i, j int) bool { return l[i].pos < l[j].pos }
func (l Reports) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

// Visits files for used nodes.
func (p *Package) Visit(node ast.Node) ast.Visitor {
	u := usedWalker(*p) // hopefully p fields are references.
	switch n := node.(type) {
	// don't walk whole file, but only:
	case *ast.ValueSpec:
		// - variable initializers
		for _, value := range n.Values {
			ast.Walk(&u, value)
		}
		// variable types.
		if n.Type != nil {
			ast.Walk(&u, n.Type)
		}
	case *ast.BlockStmt:
		// - function bodies
		for _, stmt := range n.List {
			ast.Walk(&u, stmt)
		}
	case *ast.FuncDecl:
		// - function signatures
		ast.Walk(&u, n.Type)
	case *ast.TypeSpec:
		// - type declarations
		ast.Walk(&u, n.Type)
	}
	return p
}

type usedWalker Package

// Walks through the AST marking used identifiers.
func (p *usedWalker) Visit(node ast.Node) ast.Visitor {
	// just be stupid and mark all *ast.Ident
	switch n := node.(type) {
	case *ast.Ident:
		p.prog.used[n.Name] = true
	}
	return p
}
func F2() {

}

func name(n ast.Node) string {
	switch n := n.(type) {
	case *ast.Ident:
		return n.Name
	case *ast.SelectorExpr:
		return name(n.X) + "." + n.Sel.Name
	}
	return ""
}
