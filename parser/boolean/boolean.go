package boolean

import (
	"github.com/alecthomas/participle/v2"
)

type LXBoolString string

const (
	TRUE  LXBoolString = "true"
	FALSE LXBoolString = "false"
)

type LXBool struct {
	True  *LXBoolString `parser:"@('true')"`
	False *LXBoolString `parser:"|@('false')"`
}

type LXParBool struct {
	LXBoolExpr *LXBoolExpr `parser:"'(' @@ ')'"`
}

type LXBoolExpr struct {
	LXParBool *LXParBool `parser:"@@ "`
	LXBool    *LXBool    `parser:"|@@"`
}

var LXBoolParser = participle.MustBuild[LXBoolExpr]()
