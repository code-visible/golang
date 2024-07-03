package sourcecode

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

type SourceDir struct {
	Path         string
	CountGoFiles int
}

type SourceFile struct {
	Path     string
	Name     string
	Dir      int
	AST      *ast.File
	GoSource bool
	Test     bool
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
	if err != nil {
		panic(err)
	}
	sf.AST = parsed
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
