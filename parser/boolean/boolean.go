package boolean

import (
	"github.com/alecthomas/participle/v2"
)

type LiteralString string

const (
	TRUE  LiteralString = "true"
	FALSE LiteralString = "false"
)

type UnaryOpString string

const (
	NOT_TEXT UnaryOpString = "not"
	NOT_SYMB UnaryOpString = "~"
)

type UnaryOp struct {
	Not *UnaryOpString `parser:"@('not' | '~')"`
}

type Literal struct {
	Value *LiteralString `parser:"@('true' | 'false')"`
}

type Parenthetical struct {
	BooleanExpr *BooleanExpr `parser:"'(' @@ ')'"`
}

type PrimaryExpr struct {
	Literal *Literal       `parser:"@@"`
	Paren   *Parenthetical `parser:"|@@"`
}

type UnaryExpr struct {
	Ops  []UnaryOp    `parser:"@@*"`
	Expr *PrimaryExpr `parser:"@@"`
}

type BooleanExpr struct {
	Unary *UnaryExpr `parser:"@@"`
}

var BooleanParser = participle.MustBuild[BooleanExpr](
	participle.UseLookahead(2),
)
