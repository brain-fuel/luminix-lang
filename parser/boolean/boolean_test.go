package boolean

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrueFail(t *testing.T) {
	tests := []string{
		"t rue",
		"tr ue",
		"tru e",
		"ture",
	}
	for _, test := range tests {
		_, err := BooleanParser.ParseString("", test)
		errorStr := fmt.Sprintf(
			"1:1: unexpected token \"%s\" (expected PrimaryExpr)",
			strings.Fields(test)[0],
		)
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestTrue(t *testing.T) {
	input := "true"
	expectedLit := PtrToLitString(TRUE)
	expectedPosition := Position(Position{Filename: "", Offset: 0, Line: 1, Column: 1})
	res, err := BooleanParser.ParseString("", input)
	assert.NoError(t, err)
	assert.Equal(t, expectedLit, res.Expressions[0].Bool.Unary.Expr.Lit.Value)
	assert.Equal(t, expectedPosition, res.Expressions[0].Bool.Unary.Expr.Lit.Pos)
}

func TestTrueFailAfterTrueSucceed(t *testing.T) {
	tests := []string{
		"true t rue",
		"true tr ue",
		"true tru e",
		"true ture",
	}
	expectedLit := PtrToLitString(TRUE)
	expectedPosition := Position(Position{Filename: "", Offset: 0, Line: 1, Column: 1})
	for _, test := range tests {
		res, err := BooleanParser.ParseString("", test)
		assert.Equal(t, expectedLit, res.Expressions[0].Bool.Unary.Expr.Lit.Value)
		assert.Equal(t, expectedPosition, res.Expressions[0].Bool.Unary.Expr.Lit.Pos)

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
		"true\nt rue",
		"true\ntr ue",
		"true\ntru e",
		"true\nture",
	}
	for _, test := range tests {
		_, err := BooleanParser.ParseString("", test)
		errorStr := fmt.Sprintf(
			"2:1: unexpected token \"%s\" (expected <eof>)",
			strings.Fields(test)[1],
		)
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestFalseFail(t *testing.T) {
	tests := []string{
		"f alse",
		"fa lse",
		"fal se",
		"fals e",
		"flase",
	}
	for _, test := range tests {
		_, err := BooleanParser.ParseString("", test)
		errorStr := fmt.Sprintf(
			"1:1: unexpected token \"%s\" (expected PrimaryExpr)",
			strings.Fields(test)[0],
		)
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestFalse(t *testing.T) {
	input := "false"
	expected := PtrToLitString(FALSE)
	res, err := BooleanParser.ParseString("", input)
	assert.NoError(t, err)
	assert.Equal(t, expected, res.Expressions[0].Bool.Unary.Expr.Lit.Value)
}

func TestFalseFailAfterFalseSucceed(t *testing.T) {
	tests := []string{
		"false f alse",
		"false fa lse",
		"false fal se",
		"false fals e",
		"false flase",
	}
	for _, test := range tests {
		_, err := BooleanParser.ParseString("", test)
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
		"false\nf alse",
		"false\nfa lse",
		"false\nfal se",
		"false\nfals e",
		"false\nflase",
	}
	for _, test := range tests {
		_, err := BooleanParser.ParseString("", test)
		errorStr := fmt.Sprintf(
			"2:1: unexpected token \"%s\" (expected <eof>)",
			strings.Fields(test)[1],
		)
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestParensFailWithSingleTrue(t *testing.T) {
	tests := []string{
		"(true ",
		" (true",
		"( true",
	}
	for _, test := range tests {
		_, err := BooleanParser.ParseString("", test)
		errorStr := "1:7: unexpected token \"<EOF>\" (expected \")\")"
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestParens(t *testing.T) {
	input := "(true)"
	expected := PtrToLitString(TRUE)
	res, err := BooleanParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Expressions[0].Bool.Unary.Expr.Paren)
	assert.Equal(
		t,
		expected,
		res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Lit.Value,
	)
}

func TestDoubleParens(t *testing.T) {
	input := "((true))"
	expected := PtrToLitString(TRUE)
	res, err := BooleanParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Expressions[0].Bool.Unary.Expr.Paren)
	assert.NotNil(t, res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary)
	assert.NotNil(t, res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Paren)
	assert.NotNil(
		t,
		res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Paren.BooleanExpr.Unary,
	)
	assert.NotNil(
		t,
		res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Lit,
	)
	assert.NotNil(
		t,
		res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Lit.Value,
	)
	assert.Equal(
		t,
		expected,
		res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Lit.Value,
	)
}

func TestSingleParenSetPermutations(t *testing.T) {
	tests := []string{
		// 0 spaces
		"(true)",
		// 1 space
		" (true)",
		"( true)",
		"(true )",
		"(true) ",
		// 2 spaces
		"  (true)",
		" ( true)",
		" (true )",
		" (true) ",

		" ( true)",
		"(  true)",
		"( true )",
		"( true) ",

		" (true )",
		"( true )",
		"(true  )",
		"(true ) ",

		" (true) ",
		"( true) ",
		"(true ) ",
		"(true)  ",
	}
	for _, test := range tests {
		res, err := BooleanParser.ParseString("", test)
		assert.NoError(t, err)
		assert.Equal(
			t,
			PtrToLitString(TRUE),
			res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Lit.Value,
		)
	}
}

func TestNot(t *testing.T) {
	input := "not true"
	expected0 := PtrToUnaryOpString(NOT_TEXT)
	expected1 := PtrToLitString(TRUE)
	res, err := BooleanParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expected0, res.Expressions[0].Bool.Unary.Ops[0].Not)
	assert.Equal(t, expected1, res.Expressions[0].Bool.Unary.Expr.Lit.Value)
}

func TestNotNot(t *testing.T) {
	input := "not not true"
	expected0 := PtrToUnaryOpString(NOT_TEXT)
	expected1 := PtrToUnaryOpString(NOT_TEXT)
	expected2 := PtrToLitString(TRUE)
	res, err := BooleanParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expected0, res.Expressions[0].Bool.Unary.Ops[0].Not)
	assert.Equal(t, expected1, res.Expressions[0].Bool.Unary.Ops[1].Not)
	assert.Equal(t, expected2, res.Expressions[0].Bool.Unary.Expr.Lit.Value)
}

func TestNotNotFail(t *testing.T) {
	input := "notnot true"
	_, err := BooleanParser.ParseString("", input)
	assert.EqualError(t, err, "1:1: unexpected token \"notnot\" (expected PrimaryExpr)")
}

func TestNotSymb(t *testing.T) {
	expected0 := PtrToUnaryOpString(NOT_SYMB)
	expected1 := PtrToLitString(TRUE)
	tests := []string{
		"~true",
		"~ true",
		" ~true",
		" ~ true",
		" ~  true",
	}
	for _, test := range tests {
		res, err := BooleanParser.ParseString("", test)
		assert.NoError(t, err)
		assert.Equal(
			t,
			expected0,
			res.Expressions[0].Bool.Unary.Ops[0].Not,
		)
		assert.Equal(
			t,
			expected1,
			res.Expressions[0].Bool.Unary.Expr.Lit.Value,
		)
	}
}

func PtrToLitString(s LitString) *LitString {
	return &s
}

func PtrToUnaryOpString(s UnaryOpString) *UnaryOpString {
	return &s
}
