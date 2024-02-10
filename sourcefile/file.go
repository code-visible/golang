package sourcefile

import "go/ast"

type SourceFile struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Path          string   `json:"path"`
	Pkg           string   `json:"pkg"`
	CallableCount int64    `json:"callableCount"`
	AbstractCount int64    `json:"abstractCount"`
	Deps          []string `json:"deps"`

	parsed *ast.File
	pkg    *ast.Package
}
