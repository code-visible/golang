package module

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/code-visible/golang/utils"
)

type Module struct {
	Name  string   `json:"name"`
	Path  string   `json:"path"`
	Files []string `json:"files"`

	fs   *token.FileSet
	pkgs map[string]*ast.Package
}

// initialize module
func NewModule(path string) (*Module, error) {
	// make sure the given path is a directory
	err := utils.MustDir(path)
	if err != nil {
		return nil, err
	}

	// initialize module struct
	m := &Module{
		Name:  "",
		Path:  path,
		Files: nil,
		fs:    token.NewFileSet(),
		pkgs:  nil,
	}

	// TODO: decide whether the comment should be parsed
	pkgs, err := parser.ParseDir(m.fs, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	m.pkgs = pkgs

	// exported files
	m.fs.Iterate(func(f *token.File) bool {
		m.Files = append(m.Files, f.Name())
		return true
	})

	return m, nil
}
