package golang

import (
	"path"

	"github.com/code-visible/golang/utils"
)

type Pkg struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	FullName string   `json:"fullName"`
	Path     string   `json:"path"`
	Imports  []string `json:"imports"`

	sm    *SourceMap
	cs    map[string]*Callable
	as    map[string]*Abstract
	calls []*Call
	sd    *SourceDir
	p     *Project
	imps  map[string]byte
}

func NewSourcePkg(path_ string, name string, sm *SourceMap, sd *SourceDir, p *Project) *Pkg {
	return &Pkg{
		Name:     path.Base(path_),
		FullName: name,
		Path:     path_,
		Imports:  make([]string, 0, 8),
		sm:       sm,
		cs:       make(map[string]*Callable),
		as:       make(map[string]*Abstract),
		calls:    make([]*Call, 0, 8),
		sd:       sd,
		p:        p,
		imps:     make(map[string]byte),
	}
}

func (p *Pkg) SetupID() {
	p.ID = utils.Hash(p.LookupName())
}

func (p *Pkg) LookupName() string {
	return p.FullName
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

func (p *Pkg) InjectImports() {
	for i := range p.imps {
		p.Imports = append(p.Imports, i)
	}
}

func (p *Pkg) Calls() []*Call {
	return p.calls
}

func (p *Pkg) LookupCallable(name string) *Callable {
	return p.cs[name]
}

func (p *Pkg) LookupAbstract(name string) *Abstract {
	return p.as[name]
}
