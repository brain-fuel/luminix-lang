package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"acornlang.dev/lang/lexer"
	"acornlang.dev/lang/types"
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
		_, err := FileParser.ParseString("", test)
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
	expectedPosition := types.Position(types.Position{Filename: "", Offset: 0, Line: 1, Column: 1})
	res, err := FileParser.ParseString("", input)
	assert.NoError(t, err)
	assert.Equal(t, lexer.TRUE, res.Expressions[0].Bool.Unary.Expr.Lit)
	assert.Equal(t, expectedPosition, res.Expressions[0].Bool.Unary.Expr.Pos)
}

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

func TestFalseFail(t *testing.T) {
	tests := []string{
		"f alse",
		"fa lse",
		"fal se",
		"fals e",
		"flase",
	}
	for _, test := range tests {
		_, err := FileParser.ParseString("", test)
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
	res, err := FileParser.ParseString("", input)
	assert.NoError(t, err)
	assert.Equal(t, lexer.FALSE, res.Expressions[0].Bool.Unary.Expr.Lit)
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

func TestParensFailWithSingleTrue(t *testing.T) {
	tests := []struct {
		Val         string
		ExpectedPos types.Position
	}{
		{
			"(True ",
			types.Position(types.Position{Filename: "", Offset: 1, Line: 1, Column: 2}),
		},
		{
			" (True",
			types.Position(types.Position{Filename: "", Offset: 2, Line: 1, Column: 3}),
		},
		{
			"( True",
			types.Position(types.Position{Filename: "", Offset: 2, Line: 1, Column: 3}),
		},
	}
	for _, test := range tests {
		res, err := FileParser.ParseString("", test.Val)
		assert.Equal(
			t,
			lexer.TRUE,
			res.Expressions[0].Bool.Unary.Expr.Paren.Expr.Unary.Expr.Lit,
		)
		assert.Equal(
			t,
			test.ExpectedPos,
			res.Expressions[0].Bool.Unary.Expr.Paren.Expr.Unary.Expr.Pos,
		)
		errorStr := "1:7: unexpected token \"<EOF>\" (expected \")\")"
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestParens(t *testing.T) {
	input := "(True)"
	res, err := FileParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Expressions[0].Bool.Unary.Expr.Paren)
	assert.Equal(
		t,
		lexer.TRUE,
		res.Expressions[0].Bool.Unary.Expr.Paren.Expr.Unary.Expr.Lit,
	)
}

func TestDoubleParens(t *testing.T) {
	input := "((True))"
	res, err := FileParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Expressions[0].Bool.Unary.Expr.Paren)
	assert.NotNil(t, res.Expressions[0].Bool.Unary.Expr.Paren.Expr.Unary)
	assert.NotNil(t, res.Expressions[0].Bool.Unary.Expr.Paren.Expr.Unary.Expr.Paren)
	assert.NotNil(
		t,
		res.Expressions[0].Bool.Unary.Expr.Paren.Expr.Unary.Expr.Paren.Expr.Unary,
	)
	assert.NotNil(
		t,
		res.Expressions[0].Bool.Unary.Expr.Paren.Expr.Unary.Expr.Paren.Expr.Unary.Expr.Lit,
	)
	assert.NotNil(
		t,
		res.Expressions[0].Bool.Unary.Expr.Paren.Expr.Unary.Expr.Paren.Expr.Unary.Expr.Lit,
	)
	assert.Equal(
		t,
		lexer.TRUE,
		res.Expressions[0].Bool.Unary.Expr.Paren.Expr.Unary.Expr.Paren.Expr.Unary.Expr.Lit,
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
		res, err := FileParser.ParseString("", test)
		assert.NoError(t, err)
		assert.Equal(
			t,
			lexer.TRUE,
			res.Expressions[0].Bool.Unary.Expr.Paren.Expr.Unary.Expr.Lit,
		)
	}
}

func ActAndAssertUnaryOpSuccess(t *testing.T, input string, expected string) {
	res, err := FileParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expected, res.Expressions[0].Bool.Unary.Ops[0].Op)
	assert.Equal(t, lexer.TRUE, res.Expressions[0].Bool.Unary.Expr.Lit)
}

func TestUnaryOps(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			"not True",
			lexer.NOT_TEXT,
		},
		{
			"nullify True",
			lexer.NULLIFY_TEXT,
		},
		{
			"truify True",
			lexer.TRUIFY_TEXT,
		},
		{
			"id True",
			lexer.ID_TEXT,
		},
	}
	for _, test := range testCases {
		ActAndAssertUnaryOpSuccess(t, test.input, test.expected)
	}
}

func ActAndAssertUnaryOpFail(t *testing.T, input string) {
	_, err := FileParser.ParseString("", input)
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
	res, err := FileParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, lexer.NOT_TEXT, res.Expressions[0].Bool.Unary.Ops[0].Op)
	assert.Equal(t, lexer.NOT_TEXT, res.Expressions[0].Bool.Unary.Ops[1].Op)
	assert.Equal(t, lexer.TRUE, res.Expressions[0].Bool.Unary.Expr.Lit)
}

func TestNotNotFail(t *testing.T) {
	input := "notnot true"
	_, err := FileParser.ParseString("", input)
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
		res, err := FileParser.ParseString("", test)
		assert.NoError(t, err)
		assert.Equal(
			t,
			lexer.NOT_SYMB,
			res.Expressions[0].Bool.Unary.Ops[0].Op,
		)
		assert.Equal(
			t,
			lexer.TRUE,
			res.Expressions[0].Bool.Unary.Expr.Lit,
		)
	}
}

type FailingBinopTestCase struct {
	input       string
	expectedErr error
}

func produceBinopFailingTestCases(binops ...string) []FailingBinopTestCase {
	if len(binops) == 0 {
		panic("produceBinopFailingTestCases takes at least 1 string argument")
	}
	testCases := []FailingBinopTestCase{}
	for _, binop := range binops {
		newTestCases := produceSingleBinopFailingTestCaseSet(binop)
		testCases = append(testCases, newTestCases...)
	}
	return testCases
}

func produceSingleBinopFailingTestCaseSet(binop string) []FailingBinopTestCase {
	escaped := regexp.QuoteMeta(binop)
	return []FailingBinopTestCase{
		{
			input:       binop,
			expectedErr: fmt.Errorf("1:1: unexpected token \"%s\" (expected PrimaryExpr)", escaped),
		},
		{
			input:       "True " + binop,
			expectedErr: fmt.Errorf("1:6: unexpected token \"%s\" (expected <eof>)", escaped),
		},
		{
			input:       "False " + binop,
			expectedErr: fmt.Errorf("1:7: unexpected token \"%s\" (expected <eof>)", escaped),
		},
		{
			input:       "not True " + binop,
			expectedErr: fmt.Errorf("1:10: unexpected token \"%s\" (expected <eof>)", escaped),
		},
		{
			input:       "not False " + binop,
			expectedErr: fmt.Errorf("1:11: unexpected token \"%s\" (expected <eof>)", escaped),
		},
	}
}

func TestBinopFail(t *testing.T) {
	tests := produceBinopFailingTestCases(
		lexer.AND_TEXT,
		lexer.AND_SYMB,
		lexer.NAND_TEXT,
		lexer.NAND_SYMB,
		lexer.OR_TEXT,
		lexer.OR_SYMB,
		lexer.NOR_TEXT,
		lexer.NOR_SYMB,

		lexer.XNOR_TEXT,
		lexer.IFF_TEXT,
		lexer.XNOR_SYMB,
		lexer.XOR_TEXT,
		lexer.XOR_SYMB,

		lexer.IMPLIES_TEXT,
		lexer.IMPLIES_SYMB,
		lexer.IMPLIED_BY_TEXT,
		lexer.IMPLIED_BY_SYMB,

		lexer.INHIBITS_TEXT,
		lexer.INHIBITS_SYMB,
		lexer.INHIBITED_BY_TEXT,
		lexer.INHIBITED_BY_SYMB,

		lexer.LEFT_TEXT,
		lexer.LEFT_SYMB,
		lexer.RIGHT_TEXT,
		lexer.RIGHT_SYMB,

		lexer.NOT_LEFT_TEXT,
		lexer.NOT_LEFT_SYMB,
		lexer.NOT_RIGHT_TEXT,
		lexer.NOT_RIGHT_SYMB,
	)
	for _, test := range tests {
		_, err := FileParser.ParseString("", test.input)
		assert.EqualError(t, err, test.expectedErr.Error())
	}
}

type SuccessfulBinopTestCase struct {
	input    string
	expected string
}

func produceBinopTestCases(binops ...string) []SuccessfulBinopTestCase {
	if len(binops) == 0 {
		panic("produceBinopTestCases takes at least 1 string argument")
	}
	testCases := []SuccessfulBinopTestCase{}
	for _, binop := range binops {
		newTestCases := produceSingleBinopTestCaseSet(binop)
		testCases = append(testCases, newTestCases...)
	}
	return testCases
}

func produceSingleBinopTestCaseSet(binop string) []SuccessfulBinopTestCase {
	return []SuccessfulBinopTestCase{
		{
			input:    "False " + binop + " False",
			expected: binop,
		},
		{
			input:    "False " + binop + " True",
			expected: binop,
		},
		{
			input:    "True " + binop + " False",
			expected: binop,
		},
		{
			input:    "True " + binop + " True",
			expected: binop,
		},
		{
			input:    "not False " + binop + " False",
			expected: binop,
		},
		{
			input:    "not False " + binop + " True",
			expected: binop,
		},
		{
			input:    "not True " + binop + " False",
			expected: binop,
		},
		{
			input:    "not True " + binop + " True",
			expected: binop,
		},
	}
}

func TestBinopSuccess(t *testing.T) {
	tests := produceBinopTestCases(
		lexer.AND_TEXT,
		lexer.AND_SYMB,
		lexer.NAND_TEXT,
		lexer.NAND_SYMB,
		lexer.OR_TEXT,
		lexer.OR_SYMB,
		lexer.NOR_TEXT,
		lexer.NOR_SYMB,

		lexer.XNOR_TEXT,
		lexer.IFF_TEXT,
		lexer.XNOR_SYMB,
		lexer.XOR_TEXT,
		lexer.XOR_SYMB,

		lexer.IMPLIES_TEXT,
		lexer.IMPLIES_SYMB,
		lexer.IMPLIED_BY_TEXT,
		lexer.IMPLIED_BY_SYMB,

		lexer.INHIBITS_TEXT,
		lexer.INHIBITS_SYMB,
		lexer.INHIBITED_BY_TEXT,
		lexer.INHIBITED_BY_SYMB,

		lexer.LEFT_TEXT,
		lexer.LEFT_SYMB,
		lexer.RIGHT_TEXT,
		lexer.RIGHT_SYMB,

		lexer.NOT_LEFT_TEXT,
		lexer.NOT_LEFT_SYMB,
		lexer.NOT_RIGHT_TEXT,
		lexer.NOT_RIGHT_SYMB,
	)
	for _, test := range tests {
		res, err := FileParser.ParseString("", test.input)
		actual := res.Expressions[0].Bool.Rest.Op
		assert.NoError(t, err)
		assert.Equal(t, actual, test.expected)
	}
}
