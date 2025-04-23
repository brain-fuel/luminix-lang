package boolean

import (
	"acornlang.dev/lang/types"
)

type Expr struct {
	Pos   types.Position `parser:"", json:"pos"`
	Unary *UnaryExpr     `parser:"@@"`
	Rest  *ExprRest      `parser:"(@@)?"`
}

type UnaryExpr struct {
	Pos  types.Position `parser:"", json:"pos"`
	Ops  []UnaryOp      `parser:"@@*"`
	Expr *PrimaryExpr   `parser:"@@"`
}

type UnaryOp struct {
	Pos types.Position `parser:"", json:"pos"`
	Op  string         `parser:"@UnaryOpString"`
}

type ExprRest struct {
	Pos  types.Position `parser:"", json:"pos"`
	Op   string         `parser:"@BinaryOpString"`
	Expr *Expr          `parser:"@@"`
}

type PrimaryExpr struct {
	Pos   types.Position `parser:"", json:"pos"`
	Lit   string         `parser:"@LitString"`
	Paren *ParenExpr     `parser:"|@@"`
}

type ParenExpr struct {
	Pos  types.Position `parser:"", json:"pos"`
	Expr *Expr          `parser:"'(' @@ ')'"`
}
