package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

type SourceDir struct {
	Path  string
	Files int
	Pkg   bool
}

type SourceFile struct {
	Path     string
	Name     string
	Dir      *SourceDir
	AST      *ast.File
	GoSource bool
	Test     bool
	Error    string
}

func (sf *SourceFile) Parse2AST(fset *token.FileSet) {
	sf.checkGo()
	if !sf.GoSource {
		return
	}

	// TODO: decide whether the comment should be parsed
	// TODO: handle error
	p := filepath.Join(sf.Path, sf.Name)
	parsed, err := parser.ParseFile(fset, p, nil, parser.ParseComments)
	// don't stop if current file can't be parsed
	if err != nil {
		sf.Error = err.Error()
	}
	sf.AST = parsed
	sf.Dir.Pkg = true
}

func (sf *SourceFile) checkGo() {
	p := filepath.Join(sf.Path, sf.Name)
	// TODO: should we check the content detail of the file?
	if strings.HasSuffix(p, ".go") {
		if strings.HasSuffix(p, "_test.go") {
			sf.Test = true
		}
		sf.GoSource = true
	}
}
