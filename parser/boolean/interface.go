package boolean

import (
	"errors"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type UnOp string

const (
	NOT_TEXT UnOp = "not"
	NOT_SYMB UnOp = "~"
	FALSE_OP UnOp = "nullify"
	TRUE_OP  UnOp = "truify"
	ID       UnOp = "id"
)

type BinOp string

const (
	AND_TEXT  BinOp = "and"
	AND_SYMB  BinOp = "/\\"
	NAND_TEXT BinOp = "nand"
	NAND_SYMB BinOp = "~/\\"

	OR_TEXT  BinOp = "or"
	OR_SYMB  BinOp = "\\/"
	NOR_TEXT BinOp = "nor"
	NOR_SYMB BinOp = "~\\/"

	XNOR_TEXT BinOp = "xnor"
	IFF_TEXT  BinOp = "iff"
	XNOR_SYMB BinOp = "<=>"

	XOR_TEXT BinOp = "xor"
	XOR_SYMB BinOp = "<~>"

	IMPLIES_TEXT    BinOp = "implies"
	IMPLIES_SYMB    BinOp = "=>"
	IMPLIED_BY_TEXT BinOp = "is implied by"
	IMPLIED_BY_SYMB BinOp = "<="

	INHIBITS_TEXT     BinOp = "inhibits"
	INHIBITS_SYMB     BinOp = "/=>"
	INHIBITED_BY_TEXT BinOp = "is inhibited by"
	INHIBITED_BY_SYMB BinOp = "<=/"

	LEFT_TEXT  BinOp = "left"
	LEFT_SYMB  BinOp = "<s"
	RIGHT_TEXT BinOp = "right"
	RIGHT_SYMB BinOp = "s>"

	NOT_LEFT_TEXT  BinOp = "not left"
	NOT_LEFT_SYMB  BinOp = "</"
	NOT_RIGHT_TEXT BinOp = "not right"
	NOT_RIGHT_SYMB BinOp = "/>"
)

var BoolParser = participle.MustBuild[BoolExpr](
	participle.Lexer(lexer.MustSimple([]lexer.SimpleRule{
		{
			Name:    "Keyword",
			Pattern: `\b(true|false|and|nand|or|nor|not|nullify|truify|id|implies|inhibits|left|right|xnor|iff|xor)\b`,
		},
		{Name: "MultiWordOperator", Pattern: `is implied by|is inhibited by|not left|not right`},
		{Name: "Operator", Pattern: `/\\|~/\\|\\/|~\\/|<=>|<~>|~|=>|<=|/=>|<=/|<s|s>|</|/>`},
		{Name: "Paren", Pattern: `[()]`},
		{Name: "Whitespace", Pattern: `\s+`},
	})),
	participle.Elide("Whitespace"),
)

type BoolExpr struct {
	OptUnOps []UnOp       `@("not" | "~" | "nullify" | "truify" | "id")*`
	Left     *BoolTerm    `@@`
	Rest     []*BinOpTerm `@@*`
}

type BinOpTerm struct {
	BinOp BinOp     `@("and" | "/\\" | "nand" | "~/\\" | "or" | "\\/" | "nor" | "xnor" | "iff" | "<=>" | "xor" | "<~>" | "~\\/" | "implies" | "=>" | "is implied by" | "<=" | "inhibits" | "/=>" | "is inhibited by" | "<=/" | "left" | "<s" | "right" | "s>" | "not left" | "</" | "not right" | "/>")`
	Right *BoolTerm `@@`
}

type BoolTerm struct {
	ParExp *BoolExpr `"(" @@ ")"`
	Value  *string   `| @("true" | "false")`
}

func EvalBoolExpr(expr *BoolExpr) (bool, error) {
	if expr == nil {
		return false, errors.New("invalid expression: nil")
	}

	acc, err := EvalBoolTerm(expr.Left)
	if err != nil {
		return false, err
	}

	for idx := len(expr.OptUnOps) - 1; idx >= 0; idx-- {
		switch expr.OptUnOps[idx] {
		case NOT_TEXT, NOT_SYMB:
			acc = !acc
		case FALSE_OP:
			acc = false
		case TRUE_OP:
			acc = true
		case ID:
			// No change
		}
	}

	for _, binOp := range expr.Rest {
		right, err := EvalBoolTerm(binOp.Right)
		if err != nil {
			return false, err
		}

		switch binOp.BinOp {
		case AND_TEXT, AND_SYMB:
			acc = acc && right
		case NAND_TEXT, NAND_SYMB:
			acc = !(acc && right)
		case OR_TEXT, OR_SYMB:
			acc = acc || right
		case NOR_TEXT, NOR_SYMB:
			acc = !(acc || right)
		case IMPLIES_TEXT, IMPLIES_SYMB:
			acc = !acc || right
		case IMPLIED_BY_TEXT, IMPLIED_BY_SYMB:
			acc = acc || !right
		case INHIBITS_TEXT, INHIBITS_SYMB:
			acc = !acc || !right
		case INHIBITED_BY_TEXT, INHIBITED_BY_SYMB:
			acc = !acc || !right
		case LEFT_TEXT, LEFT_SYMB:
			// no change
		case RIGHT_TEXT, RIGHT_SYMB:
			acc = right
		case NOT_LEFT_TEXT, NOT_LEFT_SYMB:
			acc = !acc
		case NOT_RIGHT_TEXT, NOT_RIGHT_SYMB:
			acc = !right
		case XNOR_TEXT, XNOR_SYMB, IFF_TEXT:
			acc = acc == right
		case XOR_TEXT, XOR_SYMB:
			acc = acc != right
		}
	}
	return acc, nil
}

func EvalBoolTerm(term *BoolTerm) (bool, error) {
	if term == nil {
		return false, errors.New("invalid term: nil expression")
	}
	if term.Value != nil {
		return *term.Value == "true", nil
	} else if term.ParExp != nil {
		return EvalBoolExpr(term.ParExp)
	}
	return false, errors.New("invalid term: missing value or subexpression")
}
