package boolean

import (
	"errors"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

const (
	FALSE_OP string = "nullify"
	TRUE_OP  string = "truify"
	ID       string = "id"
)

const (
	AND_SYMB  string = "/\\"
	NAND_TEXT string = "nand"
	NAND_SYMB string = "~/\\"

	OR_TEXT  string = "or"
	OR_SYMB  string = "\\/"
	NOR_TEXT string = "nor"
	NOR_SYMB string = "~\\/"

	XNOR_TEXT string = "xnor"
	IFF_TEXT  string = "iff"
	XNOR_SYMB string = "<=>"

	XOR_TEXT string = "xor"
	XOR_SYMB string = "<~>"

	IMPLIES_TEXT    string = "implies"
	IMPLIES_SYMB    string = "=>"
	IMPLIED_BY_TEXT string = "is implied by"
	IMPLIED_BY_SYMB string = "<="

	INHIBITS_TEXT     string = "inhibits"
	INHIBITS_SYMB     string = "/=>"
	INHIBITED_BY_TEXT string = "is inhibited by"
	INHIBITED_BY_SYMB string = "<=/"

	LEFT_TEXT  string = "left"
	LEFT_SYMB  string = "<s"
	RIGHT_TEXT string = "right"
	RIGHT_SYMB string = "s>"

	NOT_LEFT_TEXT  string = "not left"
	NOT_LEFT_SYMB  string = "</"
	NOT_RIGHT_TEXT string = "not right"
	NOT_RIGHT_SYMB string = "/>"
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
	OptUnOps []string     `@("not" | "~" | "nullify" | "truify" | "id")*`
	Left     *BoolTerm    `@@`
	Rest     []*BinOpTerm `@@*`
}

type BinOpTerm struct {
	BinOp string    `@("and" | "/\\" | "nand" | "~/\\" | "or" | "\\/" | "nor" | "xnor" | "iff" | "<=>" | "xor" | "<~>" | "~\\/" | "implies" | "=>" | "is implied by" | "<=" | "inhibits" | "/=>" | "is inhibited by" | "<=/" | "left" | "<s" | "right" | "s>" | "not left" | "</" | "not right" | "/>")`
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
		// case NOT_TEXT, NOT_SYMB:
		//	acc = !acc
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
