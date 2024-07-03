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

	// directory -> pkg
	pkgIdx map[string]int
	sm     *sourcecode.SourceMap
}

func NewProject(path string) *Project {
	p := &Project{
		Pkgs:   make([]nodes.Pkg, 0, 16),
		sm:     sourcecode.NewSourceMap(path),
		pkgIdx: make(map[string]int),
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
		p.pkgIdx[dir.Path] = len(p.Pkgs)
		pkg := nodes.NewSourcePkg(p.sm, idx)
		pkg.Path = dir.Path
		p.Pkgs = append(p.Pkgs, pkg)
	}
}

// create files from source
func (p *Project) createFiles() {
	for idx, f := range p.sm.Files() {
		pkgIdx := p.pkgIdx[f.Path]
		file := nodes.NewSourceFile(p.sm, idx, &p.Pkgs[pkgIdx])
		file.Path = f.Path
		file.Name = f.Name
		file.Pkg = pkgIdx
		p.Files = append(p.Files, file)
	}
}

// retrive the nodes
func (p *Project) retriveNodes() {
	for _, f := range p.Files {
		f.EnumerateDecls()
	}

	for _, pkg := range p.Pkgs {
		p.Callables = append(p.Callables, pkg.Callables()...)
		p.Abstracts = append(p.Abstracts, pkg.Abstracts()...)
	}
}

// retrive the calls
func (p *Project) retriveCalls() {
	for _, f := range p.Files {
		f.SearchCalls()
	}

	for _, pkg := range p.Pkgs {
		p.Calls = append(p.Calls, pkg.Calls()...)
	}
}
