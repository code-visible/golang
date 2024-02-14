package sourcepkg

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/code-visible/golang/sourcefile"
	"github.com/code-visible/golang/utils"
)

type SourcePkg struct {
	ID            string                   `json:"id"`
	Name          string                   `json:"name"`
	Path          string                   `json:"path"`
	CallableCount int64                    `json:"callableCount"`
	AbstractCount int64                    `json:"abstractCount"`
	Files         []*sourcefile.SourceFile `json:"files"`
	Deps          []string                 `json:"deps"`

	parsed    *ast.Package
	filenames []string
	files     map[string]*sourcefile.SourceFile
	fset      *token.FileSet
	module    string
}

func NewSourcePkg(module string, dir string, fset *token.FileSet) (*SourcePkg, error) {
	files, err := utils.ListGoFiles(dir, true)
	if err != nil {
		return nil, err
	}

	if dir == "." {
		dir = ""
	}

	p := &SourcePkg{
		ID:            "",
		Name:          "",
		Path:          dir,
		CallableCount: -1,
		AbstractCount: -1,
		Files:         nil,
		Deps:          nil,
		parsed:        nil,
		filenames:     files,
		files:         make(map[string]*sourcefile.SourceFile),
		fset:          fset,
		module:        module,
	}

	return p, nil
}

func (p *SourcePkg) ParseFiles() error {
	for _, f := range p.filenames {
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
			if p.Path != "" {
				p.ID = fmt.Sprintf("%s/%s", p.module, p.Path)
			} else {
				p.ID = p.module
			}
		}
		p.parsed.Files[f] = fileParsed
		sf := sourcefile.NewSourceFile(p.ID, f, fileParsed, p.fset)
		sf.EnumerateCallables()
		sf.EnumerateAbstracts()
		sf.EnumerateCallHierarchy()
		p.files[f] = sf
		p.Files = append(p.Files, sf)
	}
	return nil
}
