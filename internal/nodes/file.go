package nodes

import (
	"fmt"

	"github.com/code-visible/golang/internal/sourcecode"
)

type File struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	Pkg       *Pkg     `json:"pkg"`
	Callables []string `json:"callables"`
	Abstracts []string `json:"abstracts"`
	Calls     []string `json:"calls"`
	Deps      []string `json:"deps"`

	sm  *sourcecode.SourceMap
	idx int
}

func NewSourceFile(sm *sourcecode.SourceMap, idx int) File {
	return File{
		sm:  sm,
		idx: idx,
	}
}

func (f *File) Enumerate() {
	f.enumerateDecls()
	f.searchCalls()
}

func (f *File) enumerateDecls() {
	sf := f.sm.Files()[f.idx]
	for _, d := range sf.AST.Decls {
		fmt.Println(d)
	}
}

func (f *File) searchCalls() {}
