package sourcepkg

import "go/ast"

type SourcePkg struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Path          string   `json:"path"`
	CallableCount int64    `json:"callableCount"`
	AbstractCount int64    `json:"abstractCount"`
	Files         []string `json:"files"`
	Deps          []string `json:"deps"`

	parsed *ast.Package
}
