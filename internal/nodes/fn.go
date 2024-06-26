package nodes

import (
	"go/ast"

	"github.com/code-visible/golang/internal/parsedtypes"
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
}
