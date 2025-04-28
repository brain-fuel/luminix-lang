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
		assert.Equal(t, lexer.TRUE, res.Expressions[0].Bool.Unary.Expr.Lit)
		assert.Equal(t, expectedPosition, res.Expressions[0].Bool.Unary.Expr.Pos)

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
		assert.Equal(t, lexer.TRUE, res.Expressions[0].Bool.Unary.Expr.Lit)
		assert.Equal(t, expectedPosition, res.Expressions[0].Bool.Unary.Expr.Pos)

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
