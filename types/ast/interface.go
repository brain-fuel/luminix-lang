package ast

import (
	"acornlang.dev/lang/types"
	"acornlang.dev/lang/types/ast/boolean"
)

type File struct {
	Pos        types.Position       `parser:"", json:"pos"`
	Head       *Expr                `parser:"@@"`
	Tail       []TerminatorThenExpr `parser:"(@@)*"`
	Terminator *ExprTerminator      `parser:"(@@)?"`
	EOF        string               `parser:"EOF"`
}

type Expr struct {
	Pos  types.Position `parser:"", json:"pos"`
	Bool *boolean.Expr  `parser:"@@"`
}

type TerminatorThenExpr struct {
	Pos            types.Position  `parser:"", json:"pos"`
	ExprTerminator *ExprTerminator `parser:"@@"`
	Expr           *Expr           `parser:"@@"`
}

type ExprTerminator struct {
	Pos types.Position `parser:"", json:"pos"`
	Val []string       `parser:"@(DoubleSemicolon|Newline)+"`
}
