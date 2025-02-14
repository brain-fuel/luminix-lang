package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"

	"acornlang.dev/lang/parser/boolean"
	"acornlang.dev/lang/repl"
)

func main() {
	interactiveRepl()
}

type HistoryEntry struct {
	Raw     string
	Math    string
	English string
}

type DisplayMode int

const (
	RAW DisplayMode = iota
	MATH
	ENG
)

func interactiveRepl() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating screen: %v\n", err)
		os.Exit(1)
	}

	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing screen: %v\n", err)
		os.Exit(1)
	}
	defer screen.Fini()

	screen.Clear()
	modes := []string{"RAW", "MATH", "ENGLISH"}
	modeIndex := MATH
	rawInput := ""
	mathInput := ""
	englishInput := ""

	history := []HistoryEntry{}
	inputHistory := []string{}
	historyIndex := -1
	ctx := repl.NewReplContext()
	commandNum := 1
	scrollOffset := 0

	modeStyles := []tcell.Style{
		tcell.StyleDefault.Bold(true).Foreground(tcell.ColorWhite),
		tcell.StyleDefault.Bold(true).Foreground(tcell.ColorGreen),
		tcell.StyleDefault.Bold(true).Foreground(tcell.ColorBlue),
	}

	for {
		screen.Clear()
		width, height := screen.Size()

		modeText := fmt.Sprintf("[ %s ]", modes[modeIndex])
		drawText(
			screen, 0, 0, modeStyles[modeIndex], modeText+"  Press Ctrl+T to toggle",
			width,
		)

		drawText(screen, 0, height-1, tcell.StyleDefault.Bold(true), "Press ESC to exit", width)

		maxLines := height - 4
		totalLines := len(history)

		if totalLines > maxLines {
			if scrollOffset > totalLines-maxLines {
				scrollOffset = totalLines - maxLines
			}
			if scrollOffset < 0 {
				scrollOffset = 0
			}
		}

		yOffset := 2
		startIndex := scrollOffset
		endIndex := min(startIndex+maxLines, totalLines)
		for _, entry := range history[startIndex:endIndex] {
			var line string
			switch modeIndex {
			case RAW:
				line = entry.Raw
			case MATH:
				line = entry.Math
			case ENG:
				line = entry.English
			}
			yOffset += drawText(screen, 0, yOffset, tcell.StyleDefault, line, width)
		}

		prompt := fmt.Sprintf("lx(main):%03d:1> ", commandNum)

		var displayInput string
		switch modeIndex {
		case RAW:
			displayInput = rawInput
		case MATH:
			displayInput = mathInput
		case ENG:
			displayInput = englishInput
		}

		yOffset += drawText(
			screen,
			0,
			yOffset,
			tcell.StyleDefault.Bold(true),
			prompt+displayInput,
			width,
		)
		screen.Show()

		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				return
			case tcell.KeyEnter:
				if rawInput == "" {
					continue
				}

				fullExpression := strings.Join(wrapText(mathInput, width), " ")

				evaluated, newCtx := LXEvalPrint(fullExpression, ctx)
				ctx = newCtx

				history = append(history, HistoryEntry{
					Raw:     prompt + rawInput,
					Math:    prompt + mathInput,
					English: prompt + englishInput,
				})

				evaluatedLines := strings.Split(evaluated, "\n")
				for _, line := range evaluatedLines {
					history = append(history, HistoryEntry{Raw: line, Math: line, English: line})
				}

				inputHistory = append(inputHistory, rawInput)
				rawInput = ""
				mathInput = ""
				englishInput = ""
				historyIndex = -1
				commandNum++
				scrollOffset = max(0, len(history)-maxLines)

			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(rawInput) > 0 {
					rawInput = rawInput[:len(rawInput)-1]
					mathInput = replaceToMath(rawInput)
					englishInput = replaceToEnglish(rawInput)
				}
				historyIndex = -1

			case tcell.KeyCtrlT:
				modeIndex = DisplayMode((int(modeIndex) + 1) % len(modes))

			case tcell.KeyUp:
				if historyIndex == -1 {
					historyIndex = len(inputHistory)
				}
				if historyIndex > 0 {
					historyIndex--
					rawInput = inputHistory[historyIndex]
					mathInput = replaceToMath(rawInput)
					englishInput = replaceToEnglish(rawInput)
				}

			case tcell.KeyDown:
				if historyIndex >= 0 && historyIndex < len(inputHistory)-1 {
					historyIndex++
					rawInput = inputHistory[historyIndex]
					mathInput = replaceToMath(rawInput)
					englishInput = replaceToEnglish(rawInput)
				} else {
					historyIndex = -1
					rawInput = ""
					mathInput = ""
					englishInput = ""
				}

			default:
				if ev.Rune() != 0 {
					rawInput += string(ev.Rune())
					mathInput = replaceToMath(rawInput)
					englishInput = replaceToEnglish(rawInput)
					historyIndex = -1
				}
			}
		}
	}
}

func convertHistoryToMode(history []string, modeIndex int) []string {
	newHistory := make([]string, len(history))

	for i, line := range history {
		if modeIndex == 1 {
			newHistory[i] = replaceToMath(line)
		} else if modeIndex == 2 {
			newHistory[i] = replaceToEnglish(line)
		} else {
			newHistory[i] = line
		}
	}

	return newHistory
}

func drawText(s tcell.Screen, x, y int, style tcell.Style, text string, width int) int {
	lines := wrapText(text, width)
	for i, line := range lines {
		for j, ch := range line {
			s.SetContent(x+j, y+i, ch, nil, style)
		}
	}
	return len(lines)
}

func wrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	var lines []string
	for len(text) > width {
		splitAt := width
		for splitAt > 0 && text[splitAt] != ' ' {
			splitAt--
		}

		if splitAt == 0 {
			splitAt = width
		}

		if splitAt >= len(text) {
			break
		}

		lines = append(lines, text[:splitAt])
		text = text[splitAt:]

		if len(text) > 0 && text[0] == ' ' {
			text = text[1:]
		}
	}

	if len(text) > 0 {
		lines = append(lines, text)
	}
	return lines
}

func replaceToMath(input string) string {
	replacements := map[string]string{
		"there exists": "∃",
		"for all":      "∀",
		"and":          "/\\",
		"or":           "\\/",
		"not":          "~",
		"implies":      "=>",
		"iff":          "<=>",
	}

	acc := input
	for oldStr, newStr := range replacements {
		searchRegex := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(oldStr) + `\b`)
		acc = searchRegex.ReplaceAllString(acc, newStr)
	}
	return acc
}

func replaceToEnglish(input string) string {
	replacements := map[string]string{
		"∃":   " there exists ",
		"∀":   " for all ",
		"/\\": " and ",
		"\\/": " or ",
		"~":   " not ",
		"<=>": " iff ",
	}

	acc := input
	for oldStr, newStr := range replacements {
		searchRegex := regexp.MustCompile(regexp.QuoteMeta(oldStr))
		acc = searchRegex.ReplaceAllString(acc, newStr)
	}
	// This is for the `$n ==> result` line
	acc = regexp.MustCompile(`\b=>\b`).ReplaceAllString(acc, " implies ")
	acc = regexp.MustCompile(`\s+`).ReplaceAllString(acc, " ")
	return acc
}

func LXEvalPrint(input string, ctx *repl.ReplContext) (string, *repl.ReplContext) {
	parsed, err := boolean.BoolParser.ParseString("", input)
	if err != nil {
		return fmt.Sprintf("|  Error:\n|  illegal expression\n|  %s\n|  ^", input), ctx
	}

	result, err := boolean.EvalBoolExpr(parsed)
	if err != nil {
		return fmt.Sprintf("|  Error:\n|  %s", err.Error()), ctx
	}

	return fmt.Sprintf("$%d ==> %t", ctx.ExprNum(), result), ctx.BumpExprNum()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
