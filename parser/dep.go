package parser

import (
	"go/ast"
	"strings"

	"github.com/code-visible/golang/parser/utils"
)

type Dep struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Typ  string `json:"type"`

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
	d.ID = utils.Hash(d.Name)
}
