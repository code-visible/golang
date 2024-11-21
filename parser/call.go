package parser

import (
	"fmt"
	"go/token"

	"github.com/code-visible/golang/parser/parsedtypes"
)

type Call struct {
	ID        string `json:"id"`
	Caller    string `json:"caller"`
	Callee    string `json:"callee"`
	File      string `json:"file"`
	Typ       string `json:"typ"`
	Signature string `json:"signature"`
	Dep       string `json:"dep"`

	pos      token.Pos
	scope    string
	selector string
	typ      *parsedtypes.Type
	file     *File
}

func NewCall(pos token.Pos, scope string, selector string, typ *parsedtypes.Type) *Call {
	return &Call{
		pos:      pos,
		scope:    scope,
		selector: selector,
		typ:      typ,
	}
}

func (c *Call) Complete() {
	if c.typ != nil {
		c.Signature = fmt.Sprintf("(%s).%s()", c.typ, c.selector)
		return
	}
	if c.scope != "" {
		c.Signature = fmt.Sprintf("%s.%s()", c.scope, c.selector)
		return
	}
	c.Signature = fmt.Sprintf("%s()", c.selector)
}
