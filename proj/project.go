package proj

import (
	"github.com/code-visible/golang/internal/nodes"
	"github.com/code-visible/golang/internal/sourcecode"
)

type Project struct {
	Pkgs      []nodes.Pkg      `json:"pkgs"`
	Files     []nodes.File     `json:"files"`
	Abstracts []nodes.Abstract `json:"abstracts"`
	Callables []nodes.Callable `json:"callables"`
	// Call []call

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
	p.searchNodesAndCalls()
}

// create pkgs from source
func (p *Project) createPkgs() {
	for idx := range p.sm.Dirs() {
		p.Pkgs = append(p.Pkgs, nodes.NewSourcePkg(p.sm, idx))
	}
}

// create files from source
func (p *Project) createFiles() {
	for idx := range p.sm.Files() {
		p.Files = append(p.Files, nodes.NewSourceFile(p.sm, idx))
	}
}

// search the nodes and calls
func (p *Project) searchNodesAndCalls() {
	for _, f := range p.Files {
		f.Enumerate()
	}
}
