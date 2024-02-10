package sourcefile

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"

	"github.com/code-visible/golang/callhierarchy"
	"github.com/code-visible/golang/node"
)

type SourceFile struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Path          string   `json:"path"`
	Pkg           string   `json:"pkg"`
	CallableCount int64    `json:"callableCount"`
	AbstractCount int64    `json:"abstractCount"`
	Deps          []string `json:"deps"`

	abstracts map[string]*node.Abstract
	callables map[string]*node.Callable
	calls     []callhierarchy.Call
	parsed    *ast.File
	pkg       *ast.Package
	fset      *token.FileSet
}

func NewSourceFile(pkg string, path string, file *ast.File, fset *token.FileSet) *SourceFile {
	return &SourceFile{
		ID:            "",
		Name:          filepath.Base(path),
		Path:          filepath.Dir(path),
		Pkg:           pkg,
		CallableCount: -1,
		AbstractCount: -1,
		Deps:          nil,
		callables:     make(map[string]*node.Callable),
		parsed:        file,
		pkg:           nil,
		fset:          fset,
	}
}

func (sf *SourceFile) EnumerateCallables() {
	for _, decl := range sf.parsed.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		sf.callables[fn.Name.Name] = &node.Callable{
			ID:          "",
			Pos:         sf.fset.Position(fn.Pos()).String(),
			Name:        fn.Name.Name,
			Abstract:    "",
			Comment:     fn.Doc.Text(),
			File:        sf.Name,
			Pkg:         sf.Pkg,
			Typ:         "",
			Syscalls:    make([]string, 0),
			Parameters:  make([]string, 0),
			Results:     make([]string, 0),
			Description: "",
			Method:      false,
			Private:     false,
			Orphan:      false,
		}
	}
}

func (sf *SourceFile) EnumerateAbstracts() {
	for _, decl := range sf.parsed.Decls {
		if decl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.TypeSpec); ok {
					// TODO:
					sf.abstracts[spec.Name.String()] = &node.Abstract{
						ID:   "",
						Name: spec.Name.String(),
					}
				}
			}
		}
	}
}

func (sf *SourceFile) EnumerateCallHierarchy() {
	var trace []ast.Node
	ast.Inspect(sf.parsed, func(n ast.Node) bool {
		if n == nil {
			trace = trace[:len(trace)-1]
		} else {
			trace = append(trace, n)
		}

		if x, ok := n.(*ast.CallExpr); ok {
			call := callhierarchy.Call{
				ID:        "",
				Caller:    "universe",
				Callee:    "",
				File:      sf.Name,
				Typ:       "",
				CallerPos: -1,
				CalleePos: n.Pos(),
			}
			for i := len(trace) - 2; i >= 0; i-- {
				if fnDecl, ok := trace[i].(*ast.FuncDecl); ok {
					call.Caller = fnDecl.Name.Name
					call.CalleePos = trace[i].Pos()
					break
				}
			}

			switch fnID := x.Fun.(type) {
			case *ast.Ident:
				call.Callee = fnID.Name
			case *ast.SelectorExpr:
				if scope, ok := fnID.X.(*ast.Ident); ok {
					call.Callee = fmt.Sprintf("%s.%s", scope.Name, fnID.Sel.Name)
				} else {
					call.Callee = fmt.Sprintf("unknown.%s", fnID.Sel.Name)
				}
			default:
				panic("parse call error, not covered case")
			}

			sf.calls = append(sf.calls, call)
		}

		return true
	})
}
