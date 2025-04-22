package parser

import (
	"acornlang.dev/lang/lexer"
	"acornlang.dev/lang/types"
	"github.com/alecthomas/participle/v2"

	"acornlang.dev/lang/parser/boolean"
)

type File struct {
	Pos         types.Position `parser:"", json:"pos"`
	Expressions []Expr         `parser:"@@"`
	EOF         string         `parser:"EOF"`
}

var FileParser = participle.MustBuild[File](
	participle.Lexer(lexer.BooleanLexer),
	participle.Elide("Whitespace"),
)

type Expr struct {
	Pos            types.Position  `parser:"", json:"pos"`
	Bool           *boolean.Expr   `parser:"@@"`
	ExprTerminator *ExprTerminator `parser:"(@@)?"`
}

type ExprTerminator struct {
	Pos types.Position `parser:"", json:"pos"`
	Val []string       `parser:"@(DoubleSemicolon|Newline)+"`
}
