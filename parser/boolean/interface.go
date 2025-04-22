package boolean

import (
	"errors"
	"fmt"

	"acornlang.dev/lang/lexer"
	"acornlang.dev/lang/types"
	"github.com/alecthomas/participle/v2"
)

// Regd. Parsing

type Expr struct {
	Pos   types.Position `parser:"", json:"pos"`
	Unary *UnaryExpr     `parser:"@@"`
	Rest  *ExprRest      `parser:"(@@)?"`
}

var ExprParser = participle.MustBuild[Expr](
	participle.Lexer(lexer.BooleanLexer),
	participle.Elide("Whitespace"),
)

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

// Regd. Evaluation

type EvalResult struct {
	Pos     types.Position
	Payload bool
	Err     error
}

func errorEvalResult(pos types.Position, msg string) EvalResult {
	return EvalResult{
		Pos:     pos,
		Payload: false,
		Err:     errors.New(msg),
	}
}

func successEvalResult(pos types.Position, payload bool) EvalResult {
	return EvalResult{
		Pos:     pos,
		Payload: payload,
		Err:     nil,
	}
}

func errNotImplemented(pos types.Position, thing string) EvalResult {
	return errorEvalResult(pos, "TODO: implement evaluation of "+thing)
}

func errInvalid(pos types.Position, thing string, invalidThing any) EvalResult {
	errMsg := fmt.Errorf("invalid %s '%s'", thing, invalidThing).Error()
	return errorEvalResult(pos, errMsg)
}

func EvalPrimaryExpr(expr *PrimaryExpr) EvalResult {
	if expr == nil {
		return errInvalid(types.Position{}, "primary expression", "nil")
	}
	if expr.Paren != nil && expr.Lit != "" {
		return errInvalid(expr.Pos, "primary expression", "both Lit and Paren")
	}
	if expr.Paren != nil {
		return EvalParenExpr(expr.Paren)
	}
	switch expr.Lit {
	case lexer.TRUE:
		return successEvalResult(expr.Pos, true)
	case lexer.FALSE:
		return successEvalResult(expr.Pos, false)
	default:
		return errInvalid(expr.Pos, "boolean literal", expr.Lit)
	}
}

func EvalParenExpr(expr *ParenExpr) EvalResult {
	booleanExprRes := EvalExpr(expr.Expr)
	if booleanExprRes.Err != nil {
		return booleanExprRes
	}
	return successEvalResult(expr.Pos, booleanExprRes.Payload)
}

func EvalUnaryExpr(expr *UnaryExpr) EvalResult {
	if expr == nil {
		return errInvalid(types.Position{}, "unary expression", "nil")
	}
	exprRes := EvalPrimaryExpr(expr.Expr)
	if exprRes.Err != nil {
		return exprRes
	}
	acc := exprRes.Payload
	for idx := len(expr.Ops) - 1; idx >= 0; idx-- {
		switch expr.Ops[idx].Op {
		case lexer.NOT_TEXT, lexer.NOT_SYMB:
			acc = !acc
		case lexer.NULLIFY_TEXT:
			acc = false
		case lexer.TRUIFY_TEXT:
			acc = true
		case lexer.ID_TEXT:
			// No change
		default:
			return errInvalid(expr.Pos, "unary operator", expr.Ops[idx].Op)
		}
	}
	return successEvalResult(expr.Pos, acc)
}

func EvalExpr(expr *Expr) EvalResult {
	if expr == nil {
		return errInvalid(types.Position{}, "boolean expression", "nil")
	}
	unaryRes := EvalUnaryExpr(expr.Unary)
	return TransmogrifyUnaryResBasedOnRest(expr.Rest)(unaryRes)
}

func TransmogrifyUnaryResBasedOnRest(rest *ExprRest) func(EvalResult) EvalResult {
	if rest == nil {
		return func(unaryRes EvalResult) EvalResult {
			return unaryRes
		}
	}
	exprRes := EvalExpr(rest.Expr)
	if exprRes.Err != nil {
		return func(unaryRes EvalResult) EvalResult {
			return exprRes
		}
	}
	return func(unaryRes EvalResult) EvalResult {
		if unaryRes.Err != nil {
			return unaryRes
		}
		left := unaryRes.Payload
		right := exprRes.Payload
		var resPayload bool
		switch rest.Op {
		case lexer.AND_TEXT, lexer.AND_SYMB:
			resPayload = left && right
		case lexer.NAND_TEXT, lexer.NAND_SYMB:
			resPayload = !(left && right)
		case lexer.OR_TEXT, lexer.OR_SYMB:
			resPayload = left || right
		case lexer.NOR_TEXT, lexer.NOR_SYMB:
			resPayload = !(left || right)
		case lexer.IMPLIES_TEXT, lexer.IMPLIES_SYMB:
			resPayload = !left || right
		case lexer.IMPLIED_BY_TEXT, lexer.IMPLIED_BY_SYMB:
			resPayload = left || !right
		case lexer.INHIBITS_TEXT, lexer.INHIBITS_SYMB:
			resPayload = !left || !right
		case lexer.INHIBITED_BY_TEXT, lexer.INHIBITED_BY_SYMB:
			resPayload = !left || !right
		case lexer.LEFT_TEXT, lexer.LEFT_SYMB:
			// no change
		case lexer.RIGHT_TEXT, lexer.RIGHT_SYMB:
			resPayload = right
		case lexer.NOT_LEFT_TEXT, lexer.NOT_LEFT_SYMB:
			resPayload = !left
		case lexer.NOT_RIGHT_TEXT, lexer.NOT_RIGHT_SYMB:
			resPayload = !right
		case lexer.XNOR_TEXT, lexer.XNOR_SYMB, lexer.IFF_TEXT:
			resPayload = left == right
		case lexer.XOR_TEXT, lexer.XOR_SYMB:
			resPayload = left != right
		default:
			return errInvalid(exprRes.Pos, "binary operation", rest.Op)
		}
		return successEvalResult(unaryRes.Pos, resPayload)
	}
}
