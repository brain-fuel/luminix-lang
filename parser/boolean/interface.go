package boolean

import (
	"errors"
	"fmt"

	"github.com/alecthomas/participle/v2"
)

var BooleanExprParser = participle.MustBuild[BooleanExpr](
	participle.Lexer(BooleanLexer),
	participle.Elide("Whitespace"),
)

type ParseResult struct {
	Pos     Position
	Payload bool
	Err     error
}

func errorParseResult(pos Position, msg string) ParseResult {
	return ParseResult{
		Pos:     pos,
		Payload: false,
		Err:     errors.New(msg),
	}
}

func successParseResult(pos Position, payload bool) ParseResult {
	return ParseResult{
		Pos:     pos,
		Payload: payload,
		Err:     nil,
	}
}

func errNotImplemented(pos Position, thing string) ParseResult {
	return errorParseResult(pos, "TODO: implement evaluation of "+thing)
}

func errInvalid(pos Position, thing string, invalidThing any) ParseResult {
	errMsg := fmt.Errorf("invalid %s '%s'", thing, invalidThing).Error()
	return errorParseResult(pos, errMsg)
}

func EvalPrimaryExpr(expr *PrimaryExpr) ParseResult {
	if expr == nil {
		return errInvalid(Position{}, "primary expression", "nil")
	}
	if expr.Paren != nil && expr.Lit != "" {
		return errInvalid(expr.Pos, "primary expression", "both Lit and Paren")
	}
	if expr.Paren != nil {
		return EvalParenExpr(expr.Paren)
	}
	switch expr.Lit {
	case TRUE:
		return successParseResult(expr.Pos, true)
	case FALSE:
		return successParseResult(expr.Pos, false)
	default:
		return errInvalid(expr.Pos, "boolean literal", expr.Lit)
	}
}

func EvalParenExpr(expr *ParenExpr) ParseResult {
	booleanExprRes := EvalBooleanExpr(expr.BooleanExpr)
	if booleanExprRes.Err != nil {
		return booleanExprRes
	}
	return successParseResult(expr.Pos, booleanExprRes.Payload)
}

func EvalUnaryExpr(expr *UnaryExpr) ParseResult {
	if expr == nil {
		return errInvalid(Position{}, "unary expression", "nil")
	}
	exprRes := EvalPrimaryExpr(expr.Expr)
	if exprRes.Err != nil {
		return exprRes
	}
	acc := exprRes.Payload
	for idx := len(expr.Ops) - 1; idx >= 0; idx-- {
		switch expr.Ops[idx].Op {
		case NOT_TEXT, NOT_SYMB:
			acc = !acc
		case NULLIFY_TEXT:
			acc = false
		case TRUIFY_TEXT:
			acc = true
		case ID_TEXT:
			// No change
		default:
			return errInvalid(expr.Pos, "unary operator", expr.Ops[idx].Op)
		}
	}
	return successParseResult(expr.Pos, acc)
}

func EvalBooleanExpr(expr *BooleanExpr) ParseResult {
	if expr == nil {
		return errInvalid(Position{}, "boolean expression", "nil")
	}
	unaryRes := EvalUnaryExpr(expr.Unary)
	return TransmogrifyUnaryResBasedOnRest(expr.Rest)(unaryRes)
}

func TransmogrifyUnaryResBasedOnRest(rest *BooleanExprRest) func(ParseResult) ParseResult {
	if rest == nil {
		return func(unaryRes ParseResult) ParseResult {
			return unaryRes
		}
	}
	exprRes := EvalBooleanExpr(rest.Expr)
	if exprRes.Err != nil {
		return func(unaryRes ParseResult) ParseResult {
			return exprRes
		}
	}
	return func(unaryRes ParseResult) ParseResult {
		if unaryRes.Err != nil {
			return unaryRes
		}
		left := unaryRes.Payload
		right := exprRes.Payload
		var resPayload bool
		switch rest.Op {
		case AND_TEXT, AND_SYMB:
			resPayload = left && right
		case NAND_TEXT, NAND_SYMB:
			resPayload = !(left && right)
		case OR_TEXT, OR_SYMB:
			resPayload = left || right
		case NOR_TEXT, NOR_SYMB:
			resPayload = !(left || right)
		case IMPLIES_TEXT, IMPLIES_SYMB:
			resPayload = !left || right
		case IMPLIED_BY_TEXT, IMPLIED_BY_SYMB:
			resPayload = left || !right
		case INHIBITS_TEXT, INHIBITS_SYMB:
			resPayload = !left || !right
		case INHIBITED_BY_TEXT, INHIBITED_BY_SYMB:
			resPayload = !left || !right
		case LEFT_TEXT, LEFT_SYMB:
			// no change
		case RIGHT_TEXT, RIGHT_SYMB:
			resPayload = right
		case NOT_LEFT_TEXT, NOT_LEFT_SYMB:
			resPayload = !left
		case NOT_RIGHT_TEXT, NOT_RIGHT_SYMB:
			resPayload = !right
		case XNOR_TEXT, XNOR_SYMB, IFF_TEXT:
			resPayload = left == right
		case XOR_TEXT, XOR_SYMB:
			resPayload = left != right
		default:
			return errInvalid(exprRes.Pos, "binary operation", rest.Op)
		}
		return successParseResult(unaryRes.Pos, resPayload)
	}
}
