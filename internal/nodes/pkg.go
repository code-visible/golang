package nodes

import "github.com/code-visible/golang/internal/sourcecode"

type Pkg struct {
	ID   string `json:"id"`
	Path string `json:"path"`

	sm  *sourcecode.SourceMap
	idx int
}

func NewSourcePkg(sm *sourcecode.SourceMap, idx int) Pkg {
	return Pkg{
		sm:  sm,
		idx: idx,
	}
}
