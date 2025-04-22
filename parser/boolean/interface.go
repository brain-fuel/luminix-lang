package boolean

import (
	"errors"
	"fmt"

	"acornlang.dev/lang/lexer"
	"github.com/alecthomas/participle/v2"
)

var BooleanExprParser = participle.MustBuild[BooleanExpr](
	participle.Lexer(lexer.BooleanLexer),
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
	case lexer.TRUE:
		return successParseResult(expr.Pos, true)
	case lexer.FALSE:
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
		return successParseResult(unaryRes.Pos, resPayload)
	}
}
