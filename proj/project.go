package proj

import (
	"github.com/code-visible/golang/internal/callhierarchy"
	"github.com/code-visible/golang/internal/nodes"
	"github.com/code-visible/golang/internal/sourcecode"
)

type Project struct {
	Pkgs      []nodes.Pkg           `json:"pkgs"`
	Files     []nodes.File          `json:"files"`
	Abstracts []*nodes.Abstract     `json:"abstracts"`
	Callables []*nodes.Callable     `json:"callables"`
	Calls     []*callhierarchy.Call `json:"calls"`

	sm *sourcecode.SourceMap
}

func NewProject(path string) *Project {
	p := &Project{
		Pkgs: make([]nodes.Pkg, 0, 16),
		sm:   sourcecode.NewSourceMap(path),
	}

	return p
}

// scan the whole project to get the directories and files
func (p *Project) Initialize() {
	p.sm.Scan()
}

// parse all the files to find out the nodes we are interested at
func (p *Project) Parse() {
	p.createPkgs()
	p.createFiles()
	p.retriveNodes()
	p.retriveCalls()
}

// create pkgs from source
func (p *Project) createPkgs() {
	for idx, dir := range p.sm.Dirs() {
		pkg := nodes.NewSourcePkg(p.sm, idx)
		pkg.Path = dir.Path
		p.Pkgs = append(p.Pkgs, pkg)
	}
}

// create files from source
func (p *Project) createFiles() {
	for idx, f := range p.sm.Files() {
		file := nodes.NewSourceFile(p.sm, idx)
		file.Path = f.Path
		p.Files = append(p.Files, file)
	}
}

// retrive the nodes
func (p *Project) retriveNodes() {
	for _, f := range p.Files {
		f.EnumerateDecls()
		p.Callables = append(p.Callables, f.Callables()...)
		p.Abstracts = append(p.Abstracts, f.Abstracts()...)
	}
}

// retrive the calls
func (p *Project) retriveCalls() {
	for _, f := range p.Files {
		f.SearchCalls()
		p.Calls = append(p.Calls, f.Calls()...)
	}
}
