package boolean

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type LitString string

const (
	TRUE  LitString = "True"
	FALSE LitString = "False"
)

type TermString string

const (
	TERM_DBL_SEMICOLON    = ";;"
	TERM_NEWLINE          = "\n"
	TERM_CARRIAGE_RETURRN = "\r"
)

type UnaryOp struct {
	Pos Position       `parser:"", json:"pos"`
	Op  *UnaryOpString `parser:"@('not' | '~' | 'nullify' | 'truify' | 'id')"`
}

type BinaryOp struct {
	Pos Position        `parser:"", json:"pos"`
	Op  *BinaryOpString `@("and" | "/\\" | "nand" | "~/\\" | "or" | "\\/" | "nor" | "xnor" | "iff" | "<=>" | "xor" | "<~>" | "~\\/" | "implies" | "=>" | "is implied by" | "<=" | "inhibits" | "/=>" | "is inhibited by" | "<=/" | "left" | "<s" | "right" | "s>" | "not left" | "</" | "not right" | "/>")`
}

type Lit struct {
	Pos Position   `parser:"", json:"pos"`
	Val *LitString `parser:"@('True' | 'False')"`
}

type ParenExpr struct {
	Pos         Position     `parser:"", json:"pos"`
	BooleanExpr *BooleanExpr `parser:"'(' @@ ')'"`
}

type PrimaryExpr struct {
	Pos   Position   `parser:"", json:"pos"`
	Lit   *Lit       `parser:"@@"`
	Paren *ParenExpr `parser:"|@@"`
}

type UnaryExpr struct {
	Pos  Position     `parser:"", json:"pos"`
	Ops  []UnaryOp    `parser:"@@*"`
	Expr *PrimaryExpr `parser:"@@"`
}

type BinaryExprRest struct {
	Pos   Position   `parser:"", json:"pos"`
	Op    *BinaryOp  `parser:"@@"`
	Unary *UnaryExpr `parser:"@@"`
}

type BooleanExpr struct {
	Pos   Position         `parser:"", json:"pos"`
	Unary *UnaryExpr       `parser:"@@"`
	Rest  []BinaryExprRest `parser:"@@*"`
}

type ExprTerminator struct {
	Pos Position     `parser:"", json:"pos"`
	Val []TermString `parser:"@(';;' | '(\\r)?\\n')+"`
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

var BooleanParser = participle.MustBuild[File]()

type Position = lexer.Position
