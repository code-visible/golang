package nodes

import (
	"go/ast"
	"go/token"

	"github.com/code-visible/golang/internal/callhierarchy"
	"github.com/code-visible/golang/internal/parsedtypes"
	"github.com/code-visible/golang/internal/sourcecode"
)

type File struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Pkg  int    `json:"pkg"`
	// Callables []int  `json:"callables"`
	// Abstracts []int  `json:"abstracts"`
	// Calls     []int  `json:"calls"`
	// Deps      []int  `json:"deps"`

	sm    *sourcecode.SourceMap
	idx   int
	cs    map[string]*Callable
	as    map[string]*Abstract
	calls []*callhierarchy.Call
}

func NewSourceFile(sm *sourcecode.SourceMap, idx int) File {
	f := File{
		sm:    sm,
		idx:   idx,
		cs:    make(map[string]*Callable),
		as:    make(map[string]*Abstract),
		calls: make([]*callhierarchy.Call, 0, 8),
	}

	return f
}

func (f *File) EnumerateDecls() {
	sf := f.sm.Files()[f.idx]
	if sf.AST == nil {
		return
	}
	for _, d := range sf.AST.Decls {
		switch decl := d.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.TYPE && len(decl.Specs) > 0 {
				typSpec := decl.Specs[0].(*ast.TypeSpec)
				// ignore interface, type rename
				if strtType, ok := typSpec.Type.(*ast.StructType); ok {
					a := NewAbstract(typSpec.Name, strtType)
					a.File = f.idx
					a.Pkg = f.Pkg
					a.Pos = f.sm.FileSet().Position(a.ident.Pos()).String()
					a.Complete()
					f.as[a.Name] = a
				}
			}
		case *ast.FuncDecl:
			c := NewCallable(decl)
			c.Complete()
			c.File = f.idx
			c.Pkg = f.Pkg
			c.Pos = f.sm.FileSet().Position(c.ident.Pos()).String()
			f.cs[c.Name] = c
		}
	}
}

func (f *File) Callables() []*Callable {
	cs := make([]*Callable, 0, len(f.cs))
	for _, c := range f.cs {
		cs = append(cs, c)
	}
	return cs
}

func (f *File) Abstracts() []*Abstract {
	as := make([]*Abstract, 0, len(f.as))
	for _, a := range f.as {
		as = append(as, a)
	}
	return as
}

func (f *File) SearchCalls() {
	sf := f.sm.Files()[f.idx]
	if sf.AST == nil {
		return
	}
	ast.Inspect(sf.AST, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		// the current node is not a call expr.
		if !ok {
			return true
		}

		switch fn := call.Fun.(type) {
		case *ast.Ident:
			c := callhierarchy.NewCall(fn.Pos(), "", fn.Name, nil)
			f.calls = append(f.calls, c)
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
			c := callhierarchy.NewCall(fn.Pos(), scope, sel, typ)
			f.calls = append(f.calls, c)
			// ignore anonymous function
			// case *ast.FuncLit:
		}
		return true
	})

	for _, c := range f.calls {
		c.File = f.idx
		c.Complete()
	}
}

func (f *File) Calls() []*callhierarchy.Call {
	return f.calls
}
