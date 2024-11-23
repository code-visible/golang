package parser

import (
	"fmt"
	"go/ast"

	"github.com/code-visible/golang/parser/parsedtypes"
	"github.com/code-visible/golang/parser/utils"
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
	file   *File
}

func NewAbstract(ident *ast.Ident, strtTyp *ast.StructType, file *File) *Abstract {
	a := &Abstract{
		Name:   ident.Name,
		ident:  ident,
		fields: make(parsedtypes.Fields, 0, len(strtTyp.Fields.List)),
		file:   file,
	}

	for _, field := range strtTyp.Fields.List {
		a.fields.Parse(field)
	}

	return a
}

func (a *Abstract) SetupID() {
	a.ID = utils.Hash(a.LookupName())
}

func (a *Abstract) LookupName() string {
	return fmt.Sprintf("%s:%s", a.file.LookupName(), a.Name)
}

func (a *Abstract) Complete() {
	a.Fields = a.fields.List()
}
