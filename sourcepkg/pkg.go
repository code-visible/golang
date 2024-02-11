package sourcepkg

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/code-visible/golang/sourcefile"
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
	files  map[string]*sourcefile.SourceFile
	fset   *token.FileSet
}

func NewSourcePkg(dir string, fset *token.FileSet) (*SourcePkg, error) {
	files, err := utils.ListGoFiles(dir, true)
	if err != nil {
		return nil, err
	}

	p := &SourcePkg{
		ID:            "",
		Name:          "",
		Path:          dir,
		CallableCount: -1,
		AbstractCount: -1,
		Files:         files,
		Deps:          nil,
		parsed:        nil,
		files:         make(map[string]*sourcefile.SourceFile),
		fset:          fset,
	}

	return p, nil
}

func (p *SourcePkg) ParseFiles() error {
	for _, f := range p.Files {
		// TODO: decide whether the comment should be parsed
		fileParsed, err := parser.ParseFile(p.fset, f, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		if p.parsed == nil {
			p.parsed = &ast.Package{
				Name:  fileParsed.Name.Name,
				Files: make(map[string]*ast.File),
			}
			p.Name = fileParsed.Name.Name
		}
		p.parsed.Files[f] = fileParsed
		p.files[f] = sourcefile.NewSourceFile(p.Name, f, fileParsed, p.fset)
	}
	return nil
}
