package nodes

import (
	"go/ast"

	"github.com/code-visible/golang/internal/parsedtypes"
)

type Abstract struct {
	ID      string   `json:"id"`
	Pos     string   `json:"pos"`
	Name    string   `json:"name"`
	File    string   `json:"file"`
	Pkg     string   `json:"pkg"`
	Comment string   `json:"comment"`
	Fields  []string `json:"fields"`

	ident  *ast.Ident
	fields parsedtypes.Fields
}
