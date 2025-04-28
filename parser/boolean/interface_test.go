package boolean

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

func TestFalseFail(t *testing.T) {
	tests := []string{
		"f alse",
		"fa lse",
		"fal se",
		"fals e",
		"flase",
	}
	for _, test := range tests {
		_, err := ExprParser.ParseString("", test)
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
	expectedPosition := types.Position(types.Position{Filename: "", Offset: 0, Line: 1, Column: 1})
	res, err := ExprParser.ParseString("", input)
	assert.NoError(t, err)
	assert.Equal(t, lexer.FALSE, res.Unary.Expr.Lit)
	assert.Equal(t, expectedPosition, res.Unary.Expr.Pos)
}

func TestTrueFail(t *testing.T) {
	tests := []string{
		"t rue",
		"tr ue",
		"tru e",
		"ture",
	}
	for _, test := range tests {
		_, err := ExprParser.ParseString("", test)
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
	res, err := ExprParser.ParseString("", input)
	assert.NoError(t, err)
	assert.Equal(t, lexer.TRUE, res.Unary.Expr.Lit)
	assert.Equal(t, expectedPosition, res.Unary.Expr.Pos)
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
		res, err := ExprParser.ParseString("", test.Val)
		assert.Equal(
			t,
			lexer.TRUE,
			res.Unary.Expr.Paren.Expr.Unary.Expr.Lit,
		)
		assert.Equal(
			t,
			test.ExpectedPos,
			res.Unary.Expr.Paren.Expr.Unary.Expr.Pos,
		)
		errorStr := "1:7: unexpected token \"<EOF>\" (expected \")\")"
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestSingleParenSetPermutationsWithSingleTrue(t *testing.T) {
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
		res, err := ExprParser.ParseString("", test)
		assert.NoError(t, err)
		assert.Equal(
			t,
			lexer.TRUE,
			res.Unary.Expr.Paren.Expr.Unary.Expr.Lit,
		)
	}
}

func TestDoubleParensWithSingleTrue(t *testing.T) {
	input := "((True))"
	res, err := ExprParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, res.Unary.Expr.Paren)
	assert.NotNil(t, res.Unary.Expr.Paren.Expr.Unary)
	assert.NotNil(t, res.Unary.Expr.Paren.Expr.Unary.Expr.Paren)
	assert.NotNil(
		t,
		res.Unary.Expr.Paren.Expr.Unary.Expr.Paren.Expr.Unary,
	)
	assert.NotNil(
		t,
		res.Unary.Expr.Paren.Expr.Unary.Expr.Paren.Expr.Unary.Expr.Lit,
	)
	assert.NotNil(
		t,
		res.Unary.Expr.Paren.Expr.Unary.Expr.Paren.Expr.Unary.Expr.Lit,
	)
	assert.Equal(
		t,
		lexer.TRUE,
		res.Unary.Expr.Paren.Expr.Unary.Expr.Paren.Expr.Unary.Expr.Lit,
	)
}

func TestUnaryOpsFail(t *testing.T) {
	testCases := []string{
		"not",
		"~",
		"nullify",
		"truify",
		"id",
	}
	for _, test := range testCases {
		ActAndAssertUnaryOpFail(t, test)
	}
}

func ActAndAssertUnaryOpFail(t *testing.T, input string) {
	_, err := ExprParser.ParseString("", input)
	errorStr := fmt.Sprintf("1:%d: unexpected token \"<EOF>\" (expected PrimaryExpr)", len(input)+1)
	expectedErr := errors.New(errorStr)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestNotNotFail(t *testing.T) {
	input := "notnot true"
	_, err := ExprParser.ParseString("", input)
	assert.EqualError(t, err, "1:1: unexpected token \"notnot\" (expected PrimaryExpr)")
}

func TestNotNot(t *testing.T) {
	input := "not not True"
	res, err := ExprParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, lexer.NOT_TEXT, res.Unary.Ops[0].Op)
	assert.Equal(t, lexer.NOT_TEXT, res.Unary.Ops[1].Op)
	assert.Equal(t, lexer.TRUE, res.Unary.Expr.Lit)
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
		res, err := ExprParser.ParseString("", test)
		assert.NoError(t, err)
		assert.Equal(
			t,
			lexer.NOT_SYMB,
			res.Unary.Ops[0].Op,
		)
		assert.Equal(
			t,
			lexer.TRUE,
			res.Unary.Expr.Lit,
		)
	}
}

func TestUnaryOpSuccess(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			"not True",
			lexer.NOT_TEXT,
		},
		{
			"~True",
			lexer.NOT_SYMB,
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

func ActAndAssertUnaryOpSuccess(t *testing.T, input string, expected string) {
	res, err := ExprParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, expected, res.Unary.Ops[0].Op)
	assert.Equal(t, lexer.TRUE, res.Unary.Expr.Lit)
}

type FailingBinopTestCase struct {
	input       string
	expectedErr error
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
		_, err := ExprParser.ParseString("", test.input)
		assert.EqualError(t, err, test.expectedErr.Error())
	}
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
			input: "True " + binop,
			expectedErr: fmt.Errorf(
				"1:%d: unexpected token \"<EOF>\" (expected PrimaryExpr)",
				len(binop)+6,
			),
		},
		{
			input: "False " + binop,
			expectedErr: fmt.Errorf(
				"1:%d: unexpected token \"<EOF>\" (expected PrimaryExpr)",
				len(binop)+7,
			),
		},
		{
			input: "not True " + binop,
			expectedErr: fmt.Errorf(
				"1:%d: unexpected token \"<EOF>\" (expected PrimaryExpr)",
				len(binop)+10,
			),
		},
		{
			input: "not False " + binop,
			expectedErr: fmt.Errorf(
				"1:%d: unexpected token \"<EOF>\" (expected PrimaryExpr)",
				len(binop)+11,
			),
		},
	}
}

type SuccessfulBinopTestCase struct {
	input    string
	expected string
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
		res, err := ExprParser.ParseString("", test.input)
		actual := res.Rest.Op
		assert.NoError(t, err)
		assert.Equal(t, actual, test.expected)
	}
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
