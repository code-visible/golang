package parsedtypes

import (
	"fmt"
	"go/ast"
)

// Type represent a type of a identifier
type Type struct {
	// package name
	Pkg string
	// selector name
	Key string
	// whether it's a pointer
	Pointer bool
}

func (t *Type) String() string {
	var str string
	if t.Pkg != "" {
		str = fmt.Sprintf("%s.%s", t.Pkg, t.Key)
	} else {
		str = t.Key
	}
	if t.Pointer {
		return fmt.Sprintf("*%s", str)
	}
	return str
}

func (t *Type) Parse(expr ast.Expr) {
	*t = parseExpr(expr)
}

func parseExpr(t ast.Expr) Type {
	typ := Type{
		Pkg:     "",
		Key:     "",
		Pointer: false,
	}
	switch t_ := t.(type) {
	case *ast.SelectorExpr:
		typ.Pkg = t_.X.(*ast.Ident).Name
		typ.Key = t_.Sel.Name
		return typ
	case *ast.Ident:
		typ.Key = t_.Name
		return typ
	case *ast.StarExpr:
		typ_ := parseExpr(t_.X)
		typ_.Pointer = true
		return typ_
	case *ast.MapType:
		typKey := parseExpr(t_.Key)
		typVal := parseExpr(t_.Value)
		typ.Key = fmt.Sprintf("map[%s]%s", typKey.Key, typVal.Key)
		return typ
	case *ast.ArrayType:
		typVal := parseExpr(t_.Elt)
		typ.Key = fmt.Sprintf("[]%s", typVal.Key)
		return typ
	default:
		// function / channel ...
		// don't need to parse some specific types
		typ.Key = "unparsed_type"
		return typ
	}
}
