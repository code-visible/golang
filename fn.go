package golang

import (
	"fmt"
	"go/ast"

	"github.com/code-visible/golang/parsedtypes"
	"github.com/code-visible/golang/utils"
)

type Callable struct {
	ID         string   `json:"id"`
	Pos        string   `json:"pos"`
	Name       string   `json:"name"`
	Signature  string   `json:"signature"`
	Abstract   string   `json:"abstract"`
	File       string   `json:"file"`
	Pkg        string   `json:"pkg"`
	Comment    string   `json:"comment"`
	Parameters []string `json:"parameters"`
	Results    []string `json:"results"`
	Method     bool     `json:"method"`
	Private    bool     `json:"private"`
	Orphan     bool     `json:"orphan"`

	abstract string
	ident    *ast.Ident
	recv     parsedtypes.Field
	params   parsedtypes.Fields
	results  parsedtypes.Fields
	file     *File
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
	c.ID = utils.Hash(c.LookupName())
}

func (c *Callable) LookupName() string {
	if c.Method {
		return fmt.Sprintf("%s:(%s).%s", c.file.LookupName(), c.abstract, c.Name)
	}
	return fmt.Sprintf("%s:%s", c.file.LookupName(), c.Name)
}

func (c *Callable) Complete() {
	c.Parameters = c.params.List()
	c.Results = c.results.List()
	if c.recv.ID != nil {
		c.Method = true
		c.abstract = c.recv.Type.Key
	}
	c.Private = (c.Name[0] >= 'a' && c.Name[0] <= 'z') || c.Name[0] == '_'
	parametersStr := c.params.Format(",")
	resultsStr := c.results.Format(",")
	if c.Method {
		c.Signature = fmt.Sprintf("(%s).%s(%s) -> (%s)", c.abstract, c.Name, parametersStr, resultsStr)
	} else {
		c.Signature = fmt.Sprintf("%s(%s) -> (%s)", c.Name, parametersStr, resultsStr)
	}
}

func (c *Callable) SetupMethod() {
	if c.Method {
		abs := c.file.pkg.LookupAbstract(c.abstract)
		if abs != nil {
			c.Abstract = abs.ID
		}
	}
}
