package parser

type Pkg struct {
	ID   string `json:"id"`
	Path string `json:"path"`

	sm    *SourceMap
	idx   int
	cs    map[string]*Callable
	as    map[string]*Abstract
	calls []*Call
}

func NewSourcePkg(sm *SourceMap, idx int) Pkg {
	return Pkg{
		sm:    sm,
		idx:   idx,
		cs:    make(map[string]*Callable),
		as:    make(map[string]*Abstract),
		calls: make([]*Call, 0, 8),
	}
}

func (p *Pkg) Callables() []*Callable {
	cs := make([]*Callable, 0, len(p.cs))
	for _, c := range p.cs {
		cs = append(cs, c)
	}
	return cs
}

func (p *Pkg) Abstracts() []*Abstract {
	as := make([]*Abstract, 0, len(p.as))
	for _, a := range p.as {
		as = append(as, a)
	}
	return as
}

func (p *Pkg) Calls() []*Call {
	return p.calls
}
