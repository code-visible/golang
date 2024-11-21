package parser

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/code-visible/golang/parser/parsedtypes"
)

type File struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Pkg  string `json:"pkg"`
	// Callables []int  `json:"callables"`
	// Abstracts []int  `json:"abstracts"`
	// Calls     []int  `json:"calls"`
	// Deps      []int  `json:"deps"`

	sm   *SourceMap
	pkg  *Pkg
	deps map[string]*Dep
	sf   *SourceFile
}

func NewSourceFile(path string, name string, sm *SourceMap, sf *SourceFile, pkg *Pkg) File {
	f := File{
		Path: path,
		Name: name,
		sm:   sm,
		sf:   sf,
		pkg:  pkg,
		deps: make(map[string]*Dep),
	}

	return f
}

func (f *File) BuildDeps() {
	if f.sf.AST == nil {
		return
	}
	for _, imp := range f.sf.AST.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		d, ok := f.pkg.p.deps[importPath]
		if !ok {
			d = NewDep(importPath, imp)
			f.pkg.p.deps[importPath] = d
		}
		name := retriveImportName(imp)
		f.deps[name] = d
	}
}

func (f *File) EnumerateDecls() {
	if f.sf.AST == nil {
		return
	}
	for _, d := range f.sf.AST.Decls {
		switch decl := d.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.TYPE && len(decl.Specs) > 0 {
				typSpec := decl.Specs[0].(*ast.TypeSpec)
				// ignore interface, type rename
				if strtType, ok := typSpec.Type.(*ast.StructType); ok {
					a := NewAbstract(typSpec.Name, strtType, f)
					a.Pos = f.sm.FileSet().Position(a.ident.Pos()).String()
					a.Complete()
					f.pkg.as[a.Name] = a
				}
			}
		case *ast.FuncDecl:
			c := NewCallable(decl, f)
			c.Complete()
			c.Pos = f.sm.FileSet().Position(c.ident.Pos()).String()
			f.pkg.cs[c.Name] = c
		}
	}
}

func (f *File) SearchCalls() {
	if f.sf.AST == nil {
		return
	}
	ast.Inspect(f.sf.AST, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		// the current node is not a call expr.
		if !ok {
			return true
		}

		switch fn := call.Fun.(type) {
		case *ast.Ident:
			c := NewCall(fn.Pos(), "", fn.Name, nil)
			// callee := f.pkg.CallableDefinition(fn.Name)
			// c.Callee = callee.ID
			f.pkg.calls = append(f.pkg.calls, c)
		case *ast.SelectorExpr:
			scope := "unknown"
			sel := fn.Sel.Name
			// callee := -1
			var typ *parsedtypes.Type = nil
			if prefixIdent, ok := fn.X.(*ast.Ident); ok {
				scope = prefixIdent.Name
				if prefixIdent.Obj != nil {
					switch decl := prefixIdent.Obj.Decl.(type) {
					case *ast.AssignStmt:
						typ_ := decl.Rhs[0]
						if clt, ok := typ_.(*ast.CompositeLit); ok {
							typ = &parsedtypes.Type{}
							typ.Parse(clt.Type)
						}
					case *ast.Field:
						typ = &parsedtypes.Type{}
						typ.Parse(decl.Type)
					case *ast.ValueSpec:
						typ = &parsedtypes.Type{}
						typ.Parse(decl.Type)
					}
				}
			}
			// if typ != nil {
			// 	if typ.Pkg != "" {
			// 		if dep, ok := f.deps[typ.Pkg]; ok && !dep.std {
			// 			calleeDef := dep.pkg.CallableDefinition(sel)
			// 			if calleeDef != nil {
			// 				callee = calleeDef.File
			// 			}
			// 		}
			// 	} else {
			// 		calleeDef := f.pkg.CallableDefinition(sel)
			// 		if calleeDef != nil {
			// 			callee = calleeDef.File
			// 		}
			// 	}
			// }
			c := NewCall(fn.Pos(), scope, sel, typ)
			// c.Callee = callee
			f.pkg.calls = append(f.pkg.calls, c)
			// ignore anonymous function
			// case *ast.FuncLit:
		}
		return true
	})

	for _, c := range f.pkg.calls {
		c.Complete()
	}
}

func retriveImportName(imp *ast.ImportSpec) string {
	importPath := strings.Trim(imp.Path.Value, `"`)
	if imp.Name != nil {
		return imp.Name.Name
	}
	idx := strings.LastIndex(imp.Path.Value, "/")
	if idx < 0 {
		return importPath
	}
	return importPath[idx:]
}

// TODO: better impl
func isStd(importPath string) bool {
	return !strings.Contains(importPath, "/")
}
