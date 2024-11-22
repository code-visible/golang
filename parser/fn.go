package parser

import (
	"fmt"
	"go/ast"

	"github.com/code-visible/golang/parser/parsedtypes"
)

type Callable struct {
	ID          string   `json:"id"`
	Pos         string   `json:"pos"`
	Name        string   `json:"name"`
	Abstract    string   `json:"abstract"`
	File        string   `json:"file"`
	Pkg         string   `json:"pkg"`
	Typ         string   `json:"typ"`
	Comment     string   `json:"comment"`
	Syscalls    []string `json:"syscalls"`
	Parameters  []string `json:"parameters"`
	Results     []string `json:"results"`
	Description string   `json:"description"`
	Method      bool     `json:"method"`
	Private     bool     `json:"private"`
	Orphan      bool     `json:"orphan"`

	ident   *ast.Ident
	recv    parsedtypes.Field
	params  parsedtypes.Fields
	results parsedtypes.Fields
	file    *File
}

func NewCallable(decl *ast.FuncDecl, file *File) *Callable {
	pCnt := 0
	rCnt := 0
	if decl.Type.Params != nil {
		pCnt = len(decl.Type.Params.List)
	}
	if decl.Type.Results != nil {
		rCnt = len(decl.Type.Results.List)
	}
	c := &Callable{
		Name:    decl.Name.Name,
		ident:   decl.Name,
		params:  make(parsedtypes.Fields, 0, pCnt),
		results: make(parsedtypes.Fields, 0, rCnt),
		file:    file,
	}
	if pCnt > 0 {
		for _, pf := range decl.Type.Params.List {
			c.params.Parse(pf)
		}
	}
	if rCnt > 0 {
		for _, rf := range decl.Type.Results.List {
			c.results.Parse(rf)
		}
	}

	if decl.Recv != nil && len(decl.Recv.List) > 0 {
		var recv = make(parsedtypes.Fields, 0, 1)
		recv.Parse(decl.Recv.List[0])
		c.recv = recv[0]
	}

	return c
}

func (c *Callable) SetupID() {
	c.ID = fmt.Sprintf("%s/%s:%s", c.file.Path, c.file.Name, c.Name)
}

func (c *Callable) Complete() {
	c.Parameters = c.params.List()
	c.Results = c.results.List()
	if c.recv.ID != nil {
		c.Method = true
	}
}
