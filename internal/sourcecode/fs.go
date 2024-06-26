package sourcecode

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type SourceDir struct {
	Path         string
	CountGoFiles int
}

type SourceFile struct {
	Path     string
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
	parsed, err := parser.ParseFile(fset, sf.Path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sf.AST = parsed
}

func (sf *SourceFile) checkGo() {
	// TODO: should we check the content detail of the file?
	if strings.HasSuffix(sf.Path, ".go") {
		if strings.HasSuffix(sf.Path, "_test.go") {
			sf.Test = true
		}
		sf.GoSource = true
	}
}
