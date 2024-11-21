package parser

type Pkg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`

	sm    *SourceMap
	cs    map[string]*Callable
	as    map[string]*Abstract
	calls []*Call
	sd    *SourceDir
	p     *Project
}

func NewSourcePkg(path string, name string, sm *SourceMap, sd *SourceDir, p *Project) Pkg {
	return Pkg{
		Path:  path,
		sm:    sm,
		cs:    make(map[string]*Callable),
		as:    make(map[string]*Abstract),
		calls: make([]*Call, 0, 8),
		sd:    sd,
		p:     p,
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

func (p *Pkg) CallableDefinition(name string) *Callable {
	return p.cs[name]
}

func (p *Pkg) AbstractDefinition(name string) *Abstract {
	return p.as[name]
}
