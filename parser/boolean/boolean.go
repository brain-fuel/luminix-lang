package boolean

import (
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type WordBoundaryPolicy int

const (
	NoBoundary WordBoundaryPolicy = iota
	LeftBoundary
	RightBoundary
	BothBoundaries
)

type EscapedAndWBString struct {
	Value    string
	Boundary WordBoundaryPolicy
}

func NewEscapedAndWBString(value string, boundary WordBoundaryPolicy) EscapedAndWBString {
	return EscapedAndWBString{
		Value:    value,
		Boundary: boundary,
	}
}

func (wbs EscapedAndWBString) String() string {
	value := regexp.QuoteMeta(wbs.Value)
	switch wbs.Boundary {
	case LeftBoundary:
		return `\b` + value
	case RightBoundary:
		return value + `\b`
	case BothBoundaries:
		return `\b` + value + `\b`
	default:
		return value
	}
}

const (
	TRUE  string = "True"
	FALSE string = "False"
)

var (
	TRUE_WB  EscapedAndWBString = NewEscapedAndWBString(TRUE, BothBoundaries)
	FALSE_WB EscapedAndWBString = NewEscapedAndWBString(FALSE, BothBoundaries)
)

const (
	NOT_TEXT     string = "not"
	NOT_SYMB     string = "~"
	NULLIFY_TEXT string = "nullify"
	TRUIFY_TEXT  string = "truify"
	ID_TEXT      string = "id"
)

var (
	NOT_TEXT_WB     EscapedAndWBString = NewEscapedAndWBString(NOT_TEXT, BothBoundaries)
	NULLIFY_TEXT_WB EscapedAndWBString = NewEscapedAndWBString(NULLIFY_TEXT, BothBoundaries)
	TRUIFY_TEXT_WB  EscapedAndWBString = NewEscapedAndWBString(TRUIFY_TEXT, BothBoundaries)
	ID_TEXT_WB      EscapedAndWBString = NewEscapedAndWBString(ID_TEXT, BothBoundaries)
)

const (
	AND_TEXT  string = "and"
	AND_SYMB  string = "/\\"
	NAND_TEXT string = "nand"
	NAND_SYMB string = "~/\\"
	OR_TEXT   string = "or"
	OR_SYMB   string = "\\/"
	NOR_TEXT  string = "nor"
	NOR_SYMB  string = "~\\/"

	XNOR_TEXT string = "xnor"
	IFF_TEXT  string = "iff"
	XNOR_SYMB string = "<=>"

	XOR_TEXT string = "xor"
	XOR_SYMB string = "<~>"

	INHIBITS_TEXT     string = "inhibits"
	INHIBITS_SYMB     string = "/=>"
	INHIBITED_BY_TEXT string = "is inhibited by"
	INHIBITED_BY_SYMB string = "<=/"

	IMPLIES_TEXT    string = "implies"
	IMPLIES_SYMB    string = "=>"
	IMPLIED_BY_TEXT string = "is implied by"
	IMPLIED_BY_SYMB string = "<="

	LEFT_TEXT  string = "left"
	LEFT_SYMB  string = "<s"
	RIGHT_TEXT string = "right"
	RIGHT_SYMB string = "s>"

	NOT_LEFT_TEXT  string = "not left"
	NOT_LEFT_SYMB  string = "</"
	NOT_RIGHT_TEXT string = "not right"
	NOT_RIGHT_SYMB string = "/>"
)

var (
	AND_TEXT_WB  EscapedAndWBString = NewEscapedAndWBString(AND_TEXT, BothBoundaries)
	AND_SYMB_WB  EscapedAndWBString = NewEscapedAndWBString(AND_SYMB, NoBoundary)
	NAND_TEXT_WB EscapedAndWBString = NewEscapedAndWBString(NAND_TEXT, BothBoundaries)
	NAND_SYMB_WB EscapedAndWBString = NewEscapedAndWBString(NAND_SYMB, NoBoundary)
	OR_TEXT_WB   EscapedAndWBString = NewEscapedAndWBString(OR_TEXT, BothBoundaries)
	OR_SYMB_WB   EscapedAndWBString = NewEscapedAndWBString(OR_SYMB, NoBoundary)
	NOR_TEXT_WB  EscapedAndWBString = NewEscapedAndWBString(NOR_TEXT, BothBoundaries)
	NOR_SYMB_WB  EscapedAndWBString = NewEscapedAndWBString(NOR_SYMB, NoBoundary)

	XNOR_TEXT_WB EscapedAndWBString = NewEscapedAndWBString(XNOR_TEXT, BothBoundaries)
	IFF_TEXT_WB  EscapedAndWBString = NewEscapedAndWBString(IFF_TEXT, NoBoundary)
	XNOR_SYMB_WB EscapedAndWBString = NewEscapedAndWBString(XNOR_SYMB, BothBoundaries)

	XOR_TEXT_WB EscapedAndWBString = NewEscapedAndWBString(XOR_TEXT, BothBoundaries)
	XOR_SYMB_WB EscapedAndWBString = NewEscapedAndWBString(XOR_SYMB, NoBoundary)

	IMPLIES_TEXT_WB    EscapedAndWBString = NewEscapedAndWBString(IMPLIES_TEXT, BothBoundaries)
	IMPLIES_SYMB_WB    EscapedAndWBString = NewEscapedAndWBString(IMPLIES_SYMB, NoBoundary)
	IMPLIED_BY_TEXT_WB EscapedAndWBString = NewEscapedAndWBString(IMPLIED_BY_TEXT, BothBoundaries)
	IMPLIED_BY_SYMB_WB EscapedAndWBString = NewEscapedAndWBString(IMPLIED_BY_SYMB, NoBoundary)

	INHIBITS_TEXT_WB     EscapedAndWBString = NewEscapedAndWBString(INHIBITS_TEXT, BothBoundaries)
	INHIBITS_SYMB_WB     EscapedAndWBString = NewEscapedAndWBString(INHIBITS_SYMB, NoBoundary)
	INHIBITED_BY_TEXT_WB EscapedAndWBString = NewEscapedAndWBString(
		INHIBITED_BY_TEXT,
		BothBoundaries,
	)
	INHIBITED_BY_SYMB_WB EscapedAndWBString = NewEscapedAndWBString(INHIBITED_BY_SYMB, NoBoundary)

	LEFT_TEXT_WB  EscapedAndWBString = NewEscapedAndWBString(LEFT_TEXT, BothBoundaries)
	LEFT_SYMB_WB  EscapedAndWBString = NewEscapedAndWBString(LEFT_SYMB, NoBoundary)
	RIGHT_TEXT_WB EscapedAndWBString = NewEscapedAndWBString(RIGHT_TEXT, BothBoundaries)
	RIGHT_SYMB_WB EscapedAndWBString = NewEscapedAndWBString(RIGHT_SYMB, NoBoundary)

	NOT_LEFT_TEXT_WB  EscapedAndWBString = NewEscapedAndWBString(NOT_LEFT_TEXT, BothBoundaries)
	NOT_LEFT_SYMB_WB  EscapedAndWBString = NewEscapedAndWBString(NOT_LEFT_SYMB, NoBoundary)
	NOT_RIGHT_TEXT_WB EscapedAndWBString = NewEscapedAndWBString(NOT_RIGHT_TEXT, BothBoundaries)
	NOT_RIGHT_SYMB_WB EscapedAndWBString = NewEscapedAndWBString(NOT_RIGHT_SYMB, NoBoundary)
)

const (
	TERMINATOR_DBL_SEMICOLON    = ";;"
	TERMINATOR_NEWLINE          = "\n"
	TERMINATOR_CARRIAGE_RETURRN = "\r"
)

type TokenDef struct {
	Name   string
	Regex  string
	String string
	OneOf  []string
}

func BuildSimpleRules(tds []TokenDef) []lexer.SimpleRule {
	var rules []lexer.SimpleRule
	for _, td := range tds {
		switch {
		case td.Regex != "":
			rules = append(rules, lexer.SimpleRule{
				Name:    td.Name,
				Pattern: td.Regex,
			})
		case 0 < len(td.String):
			rules = append(rules, lexer.SimpleRule{
				Name:    td.Name,
				Pattern: td.String,
			})
		case 0 < len(td.OneOf):
			joined := "(" + strings.Join(td.OneOf, "|") + ")"
			rules = append(rules, lexer.SimpleRule{
				Name:    td.Name,
				Pattern: joined,
			})
		default:
			panic("TokenDef has no Regex or Strings: " + td.Name)
		}
	}
	return rules
}

var tokenDefinitions = []TokenDef{
	{
		Name:  "Whitespace",
		Regex: `\s+`,
	},
	{
		Name:   "DoubleSemicolon",
		String: ";;",
	},
	{
		Name:   "SingleSemicolon",
		String: ";",
	},
	{
		Name:   "LParen",
		String: "\\(",
	},
	{
		Name:   "RParen",
		String: "\\)",
	},
	{
		Name: "BinaryOpString",
		OneOf: []string{
			AND_TEXT_WB.String(),
			regexp.QuoteMeta(AND_SYMB),
			NAND_TEXT_WB.String(),
			regexp.QuoteMeta(NAND_SYMB),
			OR_TEXT_WB.String(),
			regexp.QuoteMeta(OR_SYMB),
			NOR_TEXT_WB.String(),
			regexp.QuoteMeta(NOR_SYMB),

			XNOR_TEXT_WB.String(),
			IFF_TEXT_WB.String(),
			regexp.QuoteMeta(XNOR_SYMB),

			INHIBITS_TEXT_WB.String(),
			regexp.QuoteMeta(INHIBITS_SYMB),
			INHIBITED_BY_TEXT_WB.String(),
			regexp.QuoteMeta(INHIBITED_BY_SYMB),

			XOR_TEXT_WB.String(),
			regexp.QuoteMeta(XOR_SYMB),

			IMPLIES_TEXT_WB.String(),
			regexp.QuoteMeta(IMPLIES_SYMB),
			IMPLIED_BY_TEXT_WB.String(),
			regexp.QuoteMeta(IMPLIED_BY_SYMB),

			LEFT_TEXT_WB.String(),
			regexp.QuoteMeta(LEFT_SYMB),
			RIGHT_TEXT_WB.String(),
			regexp.QuoteMeta(RIGHT_SYMB),

			NOT_LEFT_TEXT_WB.String(),
			regexp.QuoteMeta(NOT_LEFT_SYMB),
			NOT_RIGHT_TEXT_WB.String(),
			regexp.QuoteMeta(NOT_RIGHT_SYMB),
		},
	},
	{
		Name: "UnaryOpString",
		OneOf: []string{
			NOT_TEXT_WB.String(),
			`~`,
			NULLIFY_TEXT_WB.String(),
			TRUIFY_TEXT_WB.String(),
			ID_TEXT_WB.String(),
		},
	},
	{
		Name: "LitString",
		OneOf: []string{
			TRUE_WB.String(),
			FALSE_WB.String(),
		},
	},
	{
		Name:  "Newline",
		Regex: `(\r)?\n`,
	},
	{
		Name:  "Ident",
		Regex: `\b([a-zA-Z_][a-zA-Z0-9_]*)\b`,
	},
}

var simpleRules = BuildSimpleRules(tokenDefinitions)

var BooleanLexer = lexer.MustSimple(simpleRules)

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
	participle.Lexer(BooleanLexer),
	participle.Elide("Whitespace"),
)

type Position = lexer.Position
