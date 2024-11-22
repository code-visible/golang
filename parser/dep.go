package parser

import (
	"go/ast"
	"strings"
)

type Dep struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	std bool
	pkg *Pkg
}

func NewPkgDep(name string, pkg *Pkg) *Dep {
	return &Dep{
		Name: name,
		std:  false,
		pkg:  pkg,
	}
}

func NewDep(name string, imp *ast.ImportSpec) *Dep {
	importPath := strings.Trim(imp.Path.Value, `"`)
	if isStd(importPath) {
		return &Dep{
			Name: name,
			std:  true,
		}
	}
	return &Dep{
		Name: name,
		std:  false,
	}
}

func (d *Dep) SetupID() {
	d.ID = d.Name
}
