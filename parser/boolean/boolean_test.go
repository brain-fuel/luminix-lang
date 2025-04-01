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
	input := "True"
	expectedPosition := Position(Position{Filename: "", Offset: 0, Line: 1, Column: 1})
	res, err := BooleanParser.ParseString("", input)
	assert.NoError(t, err)
	assert.Equal(t, TRUE, res.Expressions[0].Bool.Unary.Expr.Lit)
	assert.Equal(t, expectedPosition, res.Expressions[0].Bool.Unary.Expr.Pos)
}

func TestTrueFailAfterTrueSucceed(t *testing.T) {
	tests := []string{
		"True t rue",
		"True tr ue",
		"True tru e",
		"True ture",
	}
	expectedPosition := Position(Position{Filename: "", Offset: 0, Line: 1, Column: 1})
	for _, test := range tests {
		res, err := BooleanParser.ParseString("", test)
		assert.Equal(t, TRUE, res.Expressions[0].Bool.Unary.Expr.Lit)
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
	expectedPosition := Position(Position{Filename: "", Offset: 0, Line: 1, Column: 1})
	for _, test := range tests {
		res, err := BooleanParser.ParseString("", test)
		assert.Equal(t, TRUE, res.Expressions[0].Bool.Unary.Expr.Lit)
		assert.Equal(t, expectedPosition, res.Expressions[0].Bool.Unary.Expr.Pos)

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
	input := "False"
	res, err := BooleanParser.ParseString("", input)
	assert.NoError(t, err)
	assert.Equal(t, FALSE, res.Expressions[0].Bool.Unary.Expr.Lit)
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
		"False\nf alse",
		"False\nfa lse",
		"False\nfal se",
		"False\nfals e",
		"False\nflase",
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
	tests := []struct {
		Val         string
		ExpectedPos Position
	}{
		{
			"(True ",
			Position(Position{Filename: "", Offset: 1, Line: 1, Column: 2}),
		},
		{
			" (True",
			Position(Position{Filename: "", Offset: 2, Line: 1, Column: 3}),
		},
		{
			"( True",
			Position(Position{Filename: "", Offset: 2, Line: 1, Column: 3}),
		},
	}
	for _, test := range tests {
		res, err := BooleanParser.ParseString("", test.Val)
		assert.Equal(
			t,
			TRUE,
			res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Lit,
		)
		assert.Equal(
			t,
			test.ExpectedPos,
			res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Pos,
		)
		errorStr := "1:7: unexpected token \"<EOF>\" (expected \")\")"
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestParens(t *testing.T) {
	input := "(True)"
	res, err := BooleanParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Expressions[0].Bool.Unary.Expr.Paren)
	assert.Equal(
		t,
		TRUE,
		res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Lit,
	)
}

func TestDoubleParens(t *testing.T) {
	input := "((True))"
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
		res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Lit,
	)
	assert.Equal(
		t,
		TRUE,
		res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Lit,
	)
}

func TestSingleParenSetPermutations(t *testing.T) {
	tests := []string{
		// 0 spaces
		"(True)",
		// 1 space
		" (True)",
		"( True)",
		"(True )",
		"(True) ",
		// 2 spaces
		"  (True)",
		" ( True)",
		" (True )",
		" (True) ",

		" ( True)",
		"(  True)",
		"( True )",
		"( True) ",

		" (True )",
		"( True )",
		"(True  )",
		"(True ) ",

		" (True) ",
		"( True) ",
		"(True ) ",
		"(True)  ",
	}
	for _, test := range tests {
		res, err := BooleanParser.ParseString("", test)
		assert.NoError(t, err)
		assert.Equal(
			t,
			TRUE,
			res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Lit,
		)
	}
}

func ActAndAssertUnaryOpSuccess(t *testing.T, input string, expected string) {
	res, err := BooleanParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expected, res.Expressions[0].Bool.Unary.Ops[0].Op)
	assert.Equal(t, TRUE, res.Expressions[0].Bool.Unary.Expr.Lit)
}

func TestUnaryOps(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			"not True",
			NOT_TEXT,
		},
		{
			"nullify True",
			NULLIFY_TEXT,
		},
		{
			"truify True",
			TRUIFY_TEXT,
		},
		{
			"id True",
			ID_TEXT,
		},
	}
	for _, test := range testCases {
		ActAndAssertUnaryOpSuccess(t, test.input, test.expected)
	}
}

func ActAndAssertUnaryOpFail(t *testing.T, input string) {
	_, err := BooleanParser.ParseString("", input)
	errorStr := fmt.Sprintf("1:%d: unexpected token \"<EOF>\" (expected PrimaryExpr)", len(input)+1)
	expectedErr := errors.New(errorStr)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestUnaryOpsFail(t *testing.T) {
	testCases := []string{
		"not",
		"nullify",
		"truify",
		"id",
	}
	for _, test := range testCases {
		ActAndAssertUnaryOpFail(t, test)
	}
}

func TestNotNot(t *testing.T) {
	input := "not not True"
	res, err := BooleanParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, NOT_TEXT, res.Expressions[0].Bool.Unary.Ops[0].Op)
	assert.Equal(t, NOT_TEXT, res.Expressions[0].Bool.Unary.Ops[1].Op)
	assert.Equal(t, TRUE, res.Expressions[0].Bool.Unary.Expr.Lit)
}

func TestNotNotFail(t *testing.T) {
	input := "notnot true"
	_, err := BooleanParser.ParseString("", input)
	assert.EqualError(t, err, "1:1: unexpected token \"notnot\" (expected PrimaryExpr)")
}

func TestNotSymb(t *testing.T) {
	tests := []string{
		"~True",
		"~ True",
		" ~True",
		" ~ True",
		" ~  True",
	}
	for _, test := range tests {
		res, err := BooleanParser.ParseString("", test)
		assert.NoError(t, err)
		assert.Equal(
			t,
			NOT_SYMB,
			res.Expressions[0].Bool.Unary.Ops[0].Op,
		)
		assert.Equal(
			t,
			TRUE,
			res.Expressions[0].Bool.Unary.Expr.Lit,
		)
	}
}
