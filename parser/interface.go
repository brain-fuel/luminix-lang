package parser

import (
	"acornlang.dev/lang/lexer"
	"acornlang.dev/lang/types/ast"
	"github.com/alecthomas/participle/v2"
)

var FileParser = participle.MustBuild[ast.File](
	participle.Lexer(lexer.BooleanLexer),
	participle.Elide("Whitespace"),
)
