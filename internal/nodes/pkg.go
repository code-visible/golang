package nodes

import "github.com/code-visible/golang/internal/sourcecode"

type Pkg struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Path  string `json:"path"`
	Files []int  `json:"files"`

	sm  *sourcecode.SourceMap
	idx int
}

func NewSourcePkg(sm *sourcecode.SourceMap, idx int) Pkg {
	return Pkg{
		sm:  sm,
		idx: idx,
	}
}
