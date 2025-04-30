package parser

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"acornlang.dev/lang/lexer"
	"acornlang.dev/lang/types"
	"github.com/stretchr/testify/assert"
)

func TestTrueFailAfterTrueSucceed(t *testing.T) {
	tests := []string{
		"True t rue",
		"True tr ue",
		"True tru e",
		"True ture",
	}
	expectedPosition := types.Position(types.Position{Filename: "", Offset: 0, Line: 1, Column: 1})
	for _, test := range tests {
		res, err := FileParser.ParseString("", test)
		assert.Equal(t, lexer.TRUE, res.Head.Bool.Unary.Expr.Lit)
		assert.Equal(t, expectedPosition, res.Head.Bool.Unary.Expr.Pos)

		errorStr := fmt.Sprintf(
			"1:6: unexpected token \"%s\" (expected <eof>)",
			strings.Fields(test)[1],
		)
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestTrueFailAfterTrueSucceedAndNewLine(t *testing.T) {
	tests := []string{
		"True\nt rue",
		"True\ntr ue",
		"True\ntru e",
		"True\nture",
	}
	expectedPosition := types.Position(types.Position{Filename: "", Offset: 0, Line: 1, Column: 1})
	for _, test := range tests {
		res, err := FileParser.ParseString("", test)
		assert.Equal(t, lexer.TRUE, res.Head.Bool.Unary.Expr.Lit)
		assert.Equal(t, expectedPosition, res.Head.Bool.Unary.Expr.Pos)

		errorStr := fmt.Sprintf(
			"2:1: unexpected token \"%s\" (expected <eof>)",
			strings.Fields(test)[1],
		)
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestFalseFailAfterFalseSucceed(t *testing.T) {
	tests := []string{
		"False f alse",
		"False fa lse",
		"False fal se",
		"False fals e",
		"False flase",
	}
	for _, test := range tests {
		_, err := FileParser.ParseString("", test)
		errorStr := fmt.Sprintf(
			"1:7: unexpected token \"%s\" (expected <eof>)",
			strings.Fields(test)[1],
		)
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestFalseFailAfterFalseSucceedAndNewLine(t *testing.T) {
	tests := []string{
		"False\nf alse",
		"False\nfa lse",
		"False\nfal se",
		"False\nfals e",
		"False\nflase",
	}
	for _, test := range tests {
		_, err := FileParser.ParseString("", test)
		errorStr := fmt.Sprintf(
			"2:1: unexpected token \"%s\" (expected <eof>)",
			strings.Fields(test)[1],
		)
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestTwoConsecutiveSuccessfulLiteralParses(t *testing.T) {
	tests := []struct {
		input          string
		expectedSecond string
	}{
		// 00
		{
			input:          "False;;False",
			expectedSecond: "False",
		},
		{
			input:          "False\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "False;;\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "False\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "False\r\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "False\n\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "False\r\n\r\nFalse",
			expectedSecond: "False",
		},
		// 01
		{
			input:          "False;;True",
			expectedSecond: "True",
		},
		{
			input:          "False\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "False;;\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "False\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "False\r\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "False\n\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "False\r\n\r\nTrue",
			expectedSecond: "True",
		},
		// 10
		{
			input:          "True;;False",
			expectedSecond: "False",
		},
		{
			input:          "True\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True;;\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True\r\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True\n\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True\r\n\r\nFalse",
			expectedSecond: "False",
		},
		// 11
		{
			input:          "True;;True",
			expectedSecond: "True",
		},
		{
			input:          "True\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True;;\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True\r\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True\n\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True\r\n\r\nTrue",
			expectedSecond: "True",
		},
	}
	for _, test := range tests {
		fmt.Println("Input:", test.input)

		res, err := FileParser.ParseString("", test.input)
		fmt.Println("Parsed Result:", res)

		assert.NoError(t, err)
		assert.Equal(t, test.expectedSecond, res.Tail[0].Expr.Bool.Unary.Expr.Lit)
	}
}

func TestOneBinaryBooleanExpressionAndThenOneBinaryLiteral(t *testing.T) {
	tests := []struct {
		input          string
		expectedSecond string
	}{
		// 00
		{
			input:          "True and False;;False",
			expectedSecond: "False",
		},
		{
			input:          "True and False\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True and False;;\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True and False\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True and False\r\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True and False\n\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True and False\r\n\r\nFalse",
			expectedSecond: "False",
		},
		// 01
		{
			input:          "True and False;;True",
			expectedSecond: "True",
		},
		{
			input:          "True and False\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True and False;;\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True and False\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True and False\r\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True and False\n\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True and False\r\n\r\nTrue",
			expectedSecond: "True",
		},
		// 10
		{
			input:          "True or False;;False",
			expectedSecond: "False",
		},
		{
			input:          "True or False\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True or False;;\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True or False\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True or False\r\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True or False\n\nFalse",
			expectedSecond: "False",
		},
		{
			input:          "True or False\r\n\r\nFalse",
			expectedSecond: "False",
		},
		// 11
		{
			input:          "True or False;;True",
			expectedSecond: "True",
		},
		{
			input:          "True or False\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True or False;;\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True or False\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True or False\r\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True or False\n\nTrue",
			expectedSecond: "True",
		},
		{
			input:          "True or False\r\n\r\nTrue",
			expectedSecond: "True",
		},
	}
	for _, test := range tests {
		res, err := FileParser.ParseString("", test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.expectedSecond, res.Tail[0].Expr.Bool.Unary.Expr.Lit)
	}
}

func TestOneBinaryLiteralAndThenOneBinaryBooleanExpression(t *testing.T) {
	tests := []struct {
		input          string
		expectedSecond string
	}{
		// 00
		{
			input:          "False;;True and False",
			expectedSecond: "False",
		},
		{
			input:          "False\nTrue and False",
			expectedSecond: "False",
		},
		{
			input:          "False;;\nTrue and False",
			expectedSecond: "False",
		},
		{
			input:          "False\nTrue and False",
			expectedSecond: "False",
		},
		{
			input:          "False\r\nTrue and False",
			expectedSecond: "False",
		},
		{
			input:          "False\n\nTrue and False",
			expectedSecond: "False",
		},
		{
			input:          "False\r\n\r\nTrue and False",
			expectedSecond: "False",
		},
		// 01
		{
			input:          "False;;False or True",
			expectedSecond: "True",
		},
		{
			input:          "False\nFalse or True",
			expectedSecond: "True",
		},
		{
			input:          "False;;\nFalse or True",
			expectedSecond: "True",
		},
		{
			input:          "False\nFalse or True",
			expectedSecond: "True",
		},
		{
			input:          "False\r\nFalse or True",
			expectedSecond: "True",
		},
		{
			input:          "False\n\nFalse or True",
			expectedSecond: "True",
		},
		{
			input:          "False\r\n\r\nFalse or True",
			expectedSecond: "True",
		},
		// 10
		{
			input:          "True;;True and False",
			expectedSecond: "False",
		},
		{
			input:          "True\nTrue and False",
			expectedSecond: "False",
		},
		{
			input:          "True;;\nTrue and False",
			expectedSecond: "False",
		},
		{
			input:          "True\nTrue and False",
			expectedSecond: "False",
		},
		{
			input:          "True\r\nTrue and False",
			expectedSecond: "False",
		},
		{
			input:          "True\n\nTrue and False",
			expectedSecond: "False",
		},
		{
			input:          "True\r\n\r\nTrue and False",
			expectedSecond: "False",
		},
		// 11
		{
			input:          "True;;False or True",
			expectedSecond: "True",
		},
		{
			input:          "True\nFalse or True",
			expectedSecond: "True",
		},
		{
			input:          "True;;\nFalse or True",
			expectedSecond: "True",
		},
		{
			input:          "True\nFalse or True",
			expectedSecond: "True",
		},
		{
			input:          "True\r\nFalse or True",
			expectedSecond: "True",
		},
		{
			input:          "True\n\nFalse or True",
			expectedSecond: "True",
		},
		{
			input:          "True\r\n\r\nFalse or True",
			expectedSecond: "True",
		},
	}
	for _, test := range tests {
		res, err := FileParser.ParseString("", test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.expectedSecond, res.Tail[0].Expr.Bool.Rest.Expr.Unary.Expr.Lit)
	}
}
