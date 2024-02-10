package callhierarchy

import (
	"go/token"
)

type Call struct {
	ID     string `json:"id"`
	Caller string `json:"caller"`
	Callee string `json:"callee"`
	File   string `json:"file"`
	Typ    string `json:"typ"`

	CallerPos token.Pos `json:"-"`
	CalleePos token.Pos `json:"-"`
}
