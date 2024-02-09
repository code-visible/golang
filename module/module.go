package module

import (
	"go/ast"
	"go/token"
)

type Module struct {
	Name  string   `json:"name"`
	Path  string   `json:"path"`
	Files []string `json:"files"`
	fs    *token.FileSet
	pkgs  map[string]*ast.Package
}
