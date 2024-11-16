package parsedtypes

import (
	"fmt"
	"go/ast"
	"strings"
)

// Field is a single field of complex structure
type Field struct {
	// the identifier of parsed AST
	ID *ast.Ident
	// type of the field
	Type Type
}

// Fields is a list of field
// Since Field is a node of a list,
// we handle fields as a whole.
type Fields []Field

func (fs *Fields) Parse(field *ast.Field) {
	t := parseExpr(field.Type)
	if len(field.Names) == 0 {
		*fs = append(*fs, Field{
			ID:   nil,
			Type: t,
		})
		return
	}
	for _, n := range field.Names {
		*fs = append(*fs, Field{
			ID:   n,
			Type: t,
		})
	}
}

func (fs *Fields) Format(seperator string) string {
	var sb strings.Builder
	for idx, p := range *fs {
		if idx > 0 {
			_, _ = sb.WriteString(seperator)
		}
		if p.ID != nil {
			_, _ = sb.WriteString(p.ID.Name)
			_, _ = sb.WriteString(" ")
			_, _ = sb.WriteString(p.Type.String())
		} else {
			_, _ = sb.WriteString(p.Type.String())
		}
	}
	return sb.String()
}

func (fs *Fields) List() []string {
	result := make([]string, 0, len(*fs))
	for _, p := range *fs {
		if p.ID != nil {
			result = append(result, fmt.Sprintf("%s %s", p.ID.Name, p.Type.String()))
		} else {
			result = append(result, p.Type.String())
		}
	}
	return result
}
