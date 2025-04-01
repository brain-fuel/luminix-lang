module acornlang.dev/lang

go 1.23.5

replace acornlang.dev/lang/parser => ./parser

replace acornlang.dev/lang/parser/boolean => ./parser/boolean

replace acornlang.dev/lang/repl => ./repl

require (
	acornlang.dev/lang/parser/boolean v0.0.0-00010101000000-000000000000
	acornlang.dev/lang/repl v0.0.0-00010101000000-000000000000
	github.com/gdamore/tcell/v2 v2.8.1
)

require (
	github.com/alecthomas/participle/v2 v2.1.4 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.3 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/term v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)
