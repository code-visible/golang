package nodes

import (
	"go/ast"

	"github.com/code-visible/golang/internal/parsedtypes"
)

type Abstract struct {
	ID      string   `json:"id"`
	Pos     string   `json:"pos"`
	Name    string   `json:"name"`
	File    int      `json:"file"`
	Pkg     int      `json:"pkg"`
	Comment string   `json:"comment"`
	Fields  []string `json:"fields"`

	ident  *ast.Ident
	fields parsedtypes.Fields
}

func NewAbstract(ident *ast.Ident, strtTyp *ast.StructType) *Abstract {
	a := &Abstract{
		Name:   ident.Name,
		ident:  ident,
		fields: make(parsedtypes.Fields, 0, len(strtTyp.Fields.List)),
	}

	for _, field := range strtTyp.Fields.List {
		a.fields.Parse(field)
	}

	return a
}

func (a *Abstract) Complete() {
	a.Fields = a.fields.List()
}
