package ast

import (
	"acornlang.dev/lang/types"
	"acornlang.dev/lang/types/ast/boolean"
)

type File struct {
	Pos         types.Position `parser:"", json:"pos"`
	Expressions []Expr         `parser:"@@"`
	EOF         string         `parser:"EOF"`
}

type Expr struct {
	Pos            types.Position  `parser:"", json:"pos"`
	Bool           *boolean.Expr   `parser:"@@"`
	ExprTerminator *ExprTerminator `parser:"(@@)?"`
}

type ExprTerminator struct {
	Pos types.Position `parser:"", json:"pos"`
	Val []string       `parser:"@(DoubleSemicolon|Newline)+"`
}
