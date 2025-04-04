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
	AND_TEXT string = "and"
	AND_SYMB string = "/\\"
)

var (
	AND_TEXT_WB EscapedAndWBString = NewEscapedAndWBString(AND_TEXT, BothBoundaries)
	AND_SYMB_WB EscapedAndWBString = NewEscapedAndWBString(AND_SYMB, NoBoundary)
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
			AND_SYMB_WB.String(),
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

var BooleanParser = participle.MustBuild[File](
	participle.Lexer(BooleanLexer),
	participle.Elide("Whitespace"),
)

type Position = lexer.Position
