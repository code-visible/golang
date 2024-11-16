package parser

type Project struct {
	Pkgs      []Pkg       `json:"pkgs"`
	Files     []File      `json:"files"`
	Abstracts []*Abstract `json:"abstracts"`
	Callables []*Callable `json:"callables"`
	Calls     []*Call     `json:"calls"`

	// directory -> pkg
	pkgIdx map[string]int
	sm     *SourceMap
}

func NewProject(gomod string, path string) *Project {
	p := &Project{
		Pkgs:   make([]Pkg, 0, 16),
		sm:     NewSourceMap(gomod, path),
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
		pkg := NewSourcePkg(p.sm, idx)
		pkg.Path = dir.Path
		p.Pkgs = append(p.Pkgs, pkg)
	}
}

// create files from source
func (p *Project) createFiles() {
	for idx, f := range p.sm.Files() {
		pkgIdx := p.pkgIdx[f.Path]
		file := NewSourceFile(p.sm, idx, &p.Pkgs[pkgIdx])
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
