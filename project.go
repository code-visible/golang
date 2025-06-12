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
	Protocol   string       `json:"protocol"`
	Timestamp  string       `json:"timestamp"`
	Repository string       `json:"repository"`
	Typ        string       `json:"typ"`
	Version    string       `json:"version"`
	Pkgs       []*Pkg       `json:"pkgs"`
	Files      []*File      `json:"files"`
	Abstracts  []*Abstract  `json:"absts"`
	Callables  []*Callable  `json:"fns"`
	Calls      []*Call      `json:"calls"`
	References []*Reference `json:"refs"`
	Deps       []*Dep       `json:"deps"`

	sm   *SourceMap
	test bool
	// pkgs    map[string]*Pkg
	dir2Pkg map[*SourceDir]*Pkg
	deps    map[string]*Dep
}

type ProjectMinify struct {
	Name       string  `json:"name"`
	Lang       string  `json:"lang"`
	Parser     string  `json:"parser"`
	Protocol   string  `json:"protocol"`
	Timestamp  string  `json:"timestamp"`
	Repository string  `json:"repository"`
	Version    string  `json:"version"`
	Typ        string  `json:"typ"`
	Abstracts  uint64  `json:"absts"`
	Callables  uint64  `json:"fns"`
	Calls      uint64  `json:"calls"`
	References uint64  `json:"refs"`
	Deps       uint64  `json:"deps"`
	Pkgs       []*Pkg  `json:"pkgs"`
	Files      []*File `json:"files"`
}

func NewProject(project, directory, excludes, module, types string, test bool) *Project {
	err := os.Chdir(project)
	if err != nil {
		panic(err)
	}
	p := &Project{
		Lang:       LANG,
		Parser:     fmt.Sprintf("%s %s", PARSER_TYPE, PARSER_VERSION),
		Protocol:   PROTOCOL_VERSION,
		Timestamp:  time.Now().Format(time.RFC3339),
		Repository: os.Getenv("repository"),
		Typ:        PARSE_TYPE_NORMAL,
		Version:    os.Getenv("version"),
		Pkgs:       make([]*Pkg, 0, 16),
		Files:      make([]*File, 0, 128),
		Abstracts:  make([]*Abstract, 0, 128),
		Callables:  make([]*Callable, 0, 1024),
		Calls:      make([]*Call, 0, 1024),
		References: make([]*Reference, 0, 128),
		Deps:       make([]*Dep, 0, 128),
		sm:         NewSourceMap(project, directory, excludes, module, types),
		test:       test,
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
		if f.GoSource {
			if !p.test && f.Test {
				continue
			}
			pkg := p.dir2Pkg[f.Dir]
			if pkg != nil {
				file := NewSourceFile(pkg.Path, f.Name, p.sm, f, pkg)
				p.Files = append(p.Files, file)
			}
			continue
		}
		file := NewSourceFile(utils.FormatPath(f.Dir.Path), f.Name, p.sm, f, nil)
		p.Files = append(p.Files, file)
	}
}

// retrive the nodes
func (p *Project) retriveNodes() {
	for _, f := range p.Files {
		if f.Source && f.Parsed {
			f.EnumerateDecls()
		}
	}

	for _, pkg := range p.Pkgs {
		p.Callables = append(p.Callables, pkg.Callables()...)
		p.Abstracts = append(p.Abstracts, pkg.Abstracts()...)
	}
}

// build dependencies of files
func (p *Project) buildDeps() {
	for _, f := range p.Files {
		if f.Source && f.Parsed {
			f.BuildDeps()
		}
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
		if f.Source && f.Parsed {
			f.SearchCalls()
		}
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
		if v.pkg != nil {
			v.Pkg = v.pkg.ID
		}
		v.SetupID()
	}

	for _, v := range p.Callables {
		v.File = v.file.ID
		v.Pkg = v.file.pkg.ID
		v.SetupID()
	}

	for _, v := range p.Abstracts {
		v.File = v.file.ID
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
		if v.Source && v.Parsed {
			v.InjectDeps()
		}
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
				if c.file != callee.file {
					c.file.Imports = append(c.file.Imports, callee.File)
				}
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
				if c.file != callee.file {
					c.file.Imports = append(c.file.Imports, callee.File)
				}
			}
		} else {
			c.Typ = CallTypePackage
		}
	}
}
