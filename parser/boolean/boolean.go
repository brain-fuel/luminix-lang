package boolean

import (
	"github.com/alecthomas/participle/v2"
)

type LXBoolString string

const (
	TRUE  LXBoolString = "true"
	FALSE LXBoolString = "false"
)

type Paren string

const (
	LPAREN Paren = "("
	RPAREN Paren = ")"
)

type LXBool struct {
	True  *LXBoolString `parser:"@('true')"`
	False *LXBoolString `parser:"|@('false')"`
}

type ParLXBool struct {
	LParen *Paren  `parser:"'('"`
	LXBool *LXBool `parser:"@@"`
	RParen *Paren  `parser:"')'"`
}

type LXBoolExpr struct {
	ParLXBool *ParLXBool `parser:"@@"`
	LXBool    LXBool     `parser:"|@@"`
}

var LXBoolParser = participle.MustBuild[LXBoolExpr]()
