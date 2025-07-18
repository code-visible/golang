package golang

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/code-visible/golang/parsedtypes"
	"github.com/code-visible/golang/utils"
)

type File struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	Source  bool     `json:"source"`
	Parsed  bool     `json:"parsed"`
	Error   string   `json:"error"`
	Pkg     string   `json:"pkg"`
	Deps    []string `json:"deps"`
	Imports []string `json:"imports"`

	sm   *SourceMap
	pkg  *Pkg
	deps map[string]*Dep
	sf   *SourceFile
	as   map[string]*Abstract
	cs   map[string]*Callable
}

func NewSourceFile(path string, name string, sm *SourceMap, sf *SourceFile, pkg *Pkg) *File {
	source := sf.GoSource
	return &File{
		Path:    path,
		Name:    name,
		Source:  source,
		Parsed:  sf.GoSource && sf.AST != nil,
		Error:   sf.Error,
		Deps:    []string{},
		Imports: make([]string, 0, 8),
		sm:      sm,
		sf:      sf,
		pkg:     pkg,
		deps:    make(map[string]*Dep),
		cs:      make(map[string]*Callable),
		as:      make(map[string]*Abstract),
	}
}

func (f *File) SetupID() {
	f.ID = utils.Hash(f.LookupName())
}

func (f *File) LookupName() string {
	if !f.Source {
		return filepath.Join("%s/%s", f.Path, f.Name)
	}
	return fmt.Sprintf("%s/%s", f.pkg.Name, f.Name)
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

func (f *File) InjectDeps() {
	for _, d := range f.deps {
		f.Deps = append(f.Deps, d.ID)
		if d.pkg != nil {
			f.pkg.imps[d.pkg.ID] = 0
		}
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
					a.Pos = filepath.ToSlash(f.sm.FileSet().Position(a.ident.Pos()).String())
					a.Complete()
					f.as[a.Name] = a
					f.pkg.as[a.Name] = a
				}
			}
		case *ast.FuncDecl:
			c := NewCallable(decl, f)
			c.Pos = filepath.ToSlash(f.sm.FileSet().Position(c.ident.Pos()).String())
			c.Complete()
			fnID := c.Name
			if c.Method {
				fnID = fmt.Sprintf("%s.%s", c.abstract, c.Name)
			}
			f.cs[fnID] = c
			f.pkg.cs[fnID] = c
		}
	}
}

func (f *File) SearchCalls() {
	if f.sf.AST == nil {
		return
	}
	var q []string
	ast.Inspect(f.sf.AST, func(n ast.Node) bool {
		callerFn, ok := n.(*ast.FuncDecl)
		if ok {
			selector := callerFn.Name.String()
			if callerFn.Recv != nil && len(callerFn.Recv.List) > 0 {
				var recv = make(parsedtypes.Fields, 0, 1)
				recv.Parse(callerFn.Recv.List[0])
				selector = fmt.Sprintf("%s.%s", recv[0].Type.Key, callerFn.Name.String())
			}
			q = append(q, selector)
			return true
		}
		call, ok := n.(*ast.CallExpr)
		// the current node is not a call expr.
		if !ok {
			return true
		}

		caller := ""
		if len(q) > 0 {
			caller = q[len(q)-1]
		}

		switch fn := call.Fun.(type) {
		case *ast.Ident:
			c := NewCall(fn.Pos(), "", fn.Name, nil, f)
			c.Pos = filepath.ToSlash(f.sm.FileSet().Position(c.pos).String())
			c.caller = caller
			f.pkg.calls = append(f.pkg.calls, c)
		case *ast.SelectorExpr:
			scope := "unknown"
			sel := fn.Sel.Name
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
			if typ != nil {
				scope = typ.Pkg
			}
			c := NewCall(fn.Pos(), scope, sel, typ, f)
			c.Pos = filepath.ToSlash(f.sm.FileSet().Position(c.pos).String())
			c.caller = caller
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

func (f *File) LookupDepByScope(scope string) *Dep {
	return f.deps[scope]
}

func (f *File) LookupCallable(name string) *Callable {
	return f.cs[name]
}

func (f *File) LookupAbstract(name string) *Abstract {
	return f.as[name]
}
