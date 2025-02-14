package boolean

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type LitString string

const (
	TRUE  LitString = "true"
	FALSE LitString = "false"
)

type UnaryOpString string

const (
	NOT_TEXT UnaryOpString = "not"
	NOT_SYMB UnaryOpString = "~"
)

type TermString string

const (
	TERM_DBL_SEMICOLON    = ";;"
	TERM_NEWLINE          = "\n"
	TERM_CARRIAGE_RETURRN = "\r"
)

type UnaryOp struct {
	Not *UnaryOpString `parser:"@('not' | '~')"`
}

type Lit struct {
	Pos   Position   `parser:"", json:"pos"`
	Value *LitString `parser:"@('true' | 'false')"`
}

type ParenExpr struct {
	BooleanExpr *BooleanExpr `parser:"'(' @@ ')'"`
}

type PrimaryExpr struct {
	Lit   *Lit       `parser:"@@"`
	Paren *ParenExpr `parser:"|@@"`
}

type UnaryExpr struct {
	Ops  []UnaryOp    `parser:"@@*"`
	Expr *PrimaryExpr `parser:"@@"`
}

type BooleanExpr struct {
	Unary *UnaryExpr `parser:"@@"`
}

type ExprTerminator struct {
	Value []TermString `parser:"@(';;' | '(\\r)?\\n')+"`
}

type Expr struct {
	Bool           *BooleanExpr    `parser:"@@"`
	ExprTerminator *ExprTerminator `parser:"(@@)?"`
}

type File struct {
	Expressions []Expr `parser:"@@"`
	EOF         string `parser:"EOF"`
}

var BooleanParser = participle.MustBuild[File]()

type Position = lexer.Position
