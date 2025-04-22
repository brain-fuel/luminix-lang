package boolean

import (
	"acornlang.dev/lang/lexer"
	"github.com/alecthomas/participle/v2"
	participleLexer "github.com/alecthomas/participle/v2/lexer"
)

type UnaryOp struct {
	Pos Position `parser:"", json:"pos"`
	Op  string   `parser:"@UnaryOpString"`
}

type ParenExpr struct {
	Pos         Position     `parser:"", json:"pos"`
	BooleanExpr *BooleanExpr `parser:"'(' @@ ')'"`
}

type PrimaryExpr struct {
	Pos   Position   `parser:"", json:"pos"`
	Lit   string     `parser:"@LitString"`
	Paren *ParenExpr `parser:"|@@"`
}

type UnaryExpr struct {
	Pos  Position     `parser:"", json:"pos"`
	Ops  []UnaryOp    `parser:"@@*"`
	Expr *PrimaryExpr `parser:"@@"`
}

type BooleanExprRest struct {
	Pos  Position     `parser:"", json:"pos"`
	Op   string       `parser:"@BinaryOpString"`
	Expr *BooleanExpr `parser:"@@"`
}

type BooleanExpr struct {
	Pos   Position         `parser:"", json:"pos"`
	Unary *UnaryExpr       `parser:"@@"`
	Rest  *BooleanExprRest `parser:"(@@)?"`
}

type ExprTerminator struct {
	Pos Position `parser:"", json:"pos"`
	Val []string `parser:"@(DoubleSemicolon|Newline)+"`
}

type Expr struct {
	Pos            Position        `parser:"", json:"pos"`
	Bool           *BooleanExpr    `parser:"@@"`
	ExprTerminator *ExprTerminator `parser:"(@@)?"`
}

type File struct {
	Pos         Position `parser:"", json:"pos"`
	Expressions []Expr   `parser:"@@"`
	EOF         string   `parser:"EOF"`
}

var FileParser = participle.MustBuild[File](
	participle.Lexer(lexer.BooleanLexer),
	participle.Elide("Whitespace"),
)

type Position = participleLexer.Position
