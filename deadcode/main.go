package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"sort"

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

	doProg(cfg.Fset, pkgs)
	return

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

func doProg(fs *token.FileSet, pkgs []*packages.Package) {
	prog := &prog{
		pkgs: pkgs,
		decl: map[string]ast.Node{},
		used: map[string]bool{},
		fs:   fs,
	}
	for _, pkg := range pkgs {
		doPackage(fs, prog, pkg)
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
		fmt.Printf("%s: %s is unused\n", fs.Position(report.pos), report.name)
	}
}
func doPackage(fs *token.FileSet, prog *prog, pkg *packages.Package) {
	p := &Package{
		p:    pkg,
		prog: prog,
	}
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
					p.prog.decl[n.Name.Name] = n
				}
			}
		}
	}
	// init() and _ are always used
	p.prog.used["init"] = true
	p.prog.used["_"] = true
	if pkg.Name != "main" {
		// exported names are marked used for non-main packages.
		// for name := range p.prog.decl {
		// 	if ast.IsExported(name) {
		// 		p.prog.used[p.p.Name+"."+name] = true
		// 	}
		// }
	} else {
		// in main programs, main() is called.
		p.prog.used["main"] = true
	}
	for _, file := range pkg.Syntax {
		// walk file looking for used nodes.
		ast.Walk(p, file)
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
