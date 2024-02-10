package sourcepkg

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/code-visible/golang/utils"
)

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

func NewSourcePkg(dir string, fset *token.FileSet) (*SourcePkg, error) {
	p := &SourcePkg{
		ID:            "",
		Name:          "",
		Path:          dir,
		CallableCount: -1,
		AbstractCount: -1,
		Files:         nil,
		Deps:          nil,
		parsed:        nil,
	}

	files, err := utils.ListGoFiles(dir, true)
	if err != nil {
		return nil, err
	}
	p.Files = files

	for _, f := range files {
		// TODO: decide whether the comment should be parsed
		fileParsed, err := parser.ParseFile(fset, f, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		if p.parsed == nil {
			p.parsed = &ast.Package{
				Name:  fileParsed.Name.Name,
				Files: make(map[string]*ast.File),
			}
			p.Name = fileParsed.Name.Name
		}
		p.parsed.Files[f] = fileParsed
	}

	return p, nil
}
