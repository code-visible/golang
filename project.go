package golang

import (
	"fmt"
	"os"
	"time"

	"github.com/code-visible/golang/utils"
)

type Project struct {
	Name       string       `json:"name"`
	Lang       string       `json:"lang"`
	Parser     string       `json:"parser"`
	Timestamp  string       `json:"timestamp"`
	Repository string       `json:"repository"`
	Version    string       `json:"version"`
	Pkgs       []*Pkg       `json:"pkgs"`
	Files      []*File      `json:"files"`
	Abstracts  []*Abstract  `json:"absts"`
	Callables  []*Callable  `json:"fns"`
	Calls      []*Call      `json:"calls"`
	References []*Reference `json:"refs"`
	Deps       []*Dep       `json:"deps"`

	sm        *SourceMap
	directory string
	// pkgs    map[string]*Pkg
	dir2Pkg map[*SourceDir]*Pkg
	deps    map[string]*Dep
}

func NewProject(project, directory string) *Project {
	err := os.Chdir(project)
	if err != nil {
		panic(err)
	}
	p := &Project{
		Lang:       LANG,
		Parser:     fmt.Sprintf("%s %s", PARSER_TYPE, VERSION),
		directory:  directory,
		Timestamp:  time.Now().Format(time.RFC3339),
		Repository: os.Getenv("repository"),
		Version:    os.Getenv("version"),
		Pkgs:       make([]*Pkg, 0, 16),
		Files:      make([]*File, 0, 128),
		Abstracts:  make([]*Abstract, 0, 128),
		Callables:  make([]*Callable, 0, 1024),
		Calls:      make([]*Call, 0, 1024),
		References: make([]*Reference, 0, 128),
		Deps:       make([]*Dep, 0, 128),
		sm:         NewSourceMap(project, directory),
		// pkgs:      make(map[string]*Pkg),
		dir2Pkg: make(map[*SourceDir]*Pkg),
		deps:    make(map[string]*Dep),
	}

	p.Name = p.sm.module

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
	p.buildDeps()
	p.retriveCalls()
	p.injectFields()
	p.connect()
}

// create pkgs from source
func (p *Project) createPkgs() {
	for _, dir := range p.sm.Dirs() {
		if !dir.Pkg {
			continue
		}
		lookupName := fmt.Sprintf("%s/%s", p.Name, dir.Path)
		if dir.Path == "." {
			lookupName = p.Name
		}
		pkg := NewSourcePkg(utils.FormatPath(dir.Path), lookupName, p.sm, dir, p)
		p.Pkgs = append(p.Pkgs, pkg)
		p.deps[lookupName] = NewPkgDep(lookupName, pkg)
		p.dir2Pkg[dir] = pkg
	}
}

// create files from source
func (p *Project) createFiles() {
	for _, f := range p.sm.Files() {
		if !f.GoSource || f.Test {
			continue
		}
		pkg := p.dir2Pkg[f.Dir]
		if pkg != nil {
			file := NewSourceFile(pkg.Path, f.Name, p.sm, f, pkg)
			p.Files = append(p.Files, file)
		}
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

// build dependencies of files
func (p *Project) buildDeps() {
	for _, f := range p.Files {
		f.BuildDeps()
	}

	for _, v := range p.deps {
		if v.std {
			v.Typ = "std"
		} else if v.pkg != nil {
			v.Typ = "pkg"
		} else {
			v.Typ = "open"
		}
		p.Deps = append(p.Deps, v)
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

func (p *Project) injectFields() {
	for _, v := range p.Pkgs {
		v.SetupID()
	}

	for _, v := range p.Files {
		v.Pkg = v.pkg.ID
		v.SetupID()
	}

	for _, v := range p.Callables {
		v.File = v.file.ID
		v.Pkg = v.file.pkg.ID
		v.SetupID()
	}

	for _, v := range p.Abstracts {
		v.File = v.file.ID
		v.Pkg = v.file.pkg.ID
		v.SetupID()
	}

	for _, v := range p.Calls {
		v.File = v.file.ID
		v.SetupID()
	}

	for _, v := range p.Deps {
		v.SetupID()
		v.SetupRef()
	}

	for _, v := range p.Files {
		v.InjectDeps()
	}

	for _, v := range p.Callables {
		v.SetupMethod()
	}

	for _, v := range p.Pkgs {
		v.InjectImports()
	}
}

func (p *Project) connect() {
	for _, c := range p.Calls {
		caller := c.file.LookupCallable(c.caller)
		if caller != nil {
			c.Caller = caller.ID
		}
		selector := c.selector
		if c.typ != nil {
			selector = fmt.Sprintf("%s.%s", c.typ.Key, c.selector)
		}
		if c.scope == "" {
			if _, ok := Builtins[selector]; ok {
				c.Typ = CallTypeBuiltin
				continue
			}
			c.Typ = CallTypeInternal
			callee := c.file.pkg.LookupCallable(selector)
			if callee != nil {
				c.Callee = callee.ID
				c.file.Imports = append(c.file.Imports, callee.File)
			}
			continue
		}
		if _, ok := Libs[c.scope]; ok {
			c.Typ = "std"
			continue
		}
		dep := c.file.LookupDepByScope(c.scope)
		if dep != nil {
			c.Typ = CallTypeExternal
			c.Dep = dep.Name
			if dep.std || dep.pkg == nil {
				continue
			}
			callee := dep.pkg.LookupCallable(selector)
			if callee != nil {
				c.Callee = callee.ID
				c.file.Imports = append(c.file.Imports, callee.File)
			}
		} else {
			c.Typ = CallTypePackage
		}
	}
}
