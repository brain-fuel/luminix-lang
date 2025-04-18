package boolean

import (
	"errors"
	"fmt"
	"regexp"
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
	expectedPosition := Position(Position{Filename: "", Offset: 0, Line: 1, Column: 1})
	res, err := FileParser.ParseString("", input)
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
		res, err := FileParser.ParseString("", test)
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
		res, err := FileParser.ParseString("", test)
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
		res, err := FileParser.ParseString("", test.Val)
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
	res, err := FileParser.ParseString("", input)
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
	res, err := FileParser.ParseString("", input)
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
		res, err := FileParser.ParseString("", test)
		assert.NoError(t, err)
		assert.Equal(
			t,
			TRUE,
			res.Expressions[0].Bool.Unary.Expr.Paren.BooleanExpr.Unary.Expr.Lit,
		)
	}
}

func ActAndAssertUnaryOpSuccess(t *testing.T, input string, expected string) {
	res, err := FileParser.ParseString("", input)
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
	assert.Equal(t, NOT_TEXT, res.Expressions[0].Bool.Unary.Ops[0].Op)
	assert.Equal(t, NOT_TEXT, res.Expressions[0].Bool.Unary.Ops[1].Op)
	assert.Equal(t, TRUE, res.Expressions[0].Bool.Unary.Expr.Lit)
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
		AND_TEXT,
		AND_SYMB,
		NAND_TEXT,
		NAND_SYMB,
		OR_TEXT,
		OR_SYMB,
		NOR_TEXT,
		NOR_SYMB,

		XNOR_TEXT,
		IFF_TEXT,
		XNOR_SYMB,
		XOR_TEXT,
		XOR_SYMB,

		IMPLIES_TEXT,
		IMPLIES_SYMB,
		IMPLIED_BY_TEXT,
		IMPLIED_BY_SYMB,

		INHIBITS_TEXT,
		INHIBITS_SYMB,
		INHIBITED_BY_TEXT,
		INHIBITED_BY_SYMB,

		LEFT_TEXT,
		LEFT_SYMB,
		RIGHT_TEXT,
		RIGHT_SYMB,

		NOT_LEFT_TEXT,
		NOT_LEFT_SYMB,
		NOT_RIGHT_TEXT,
		NOT_RIGHT_SYMB,
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
		AND_TEXT,
		AND_SYMB,
		NAND_TEXT,
		NAND_SYMB,
		OR_TEXT,
		OR_SYMB,
		NOR_TEXT,
		NOR_SYMB,

		XNOR_TEXT,
		IFF_TEXT,
		XNOR_SYMB,
		XOR_TEXT,
		XOR_SYMB,

		IMPLIES_TEXT,
		IMPLIES_SYMB,
		IMPLIED_BY_TEXT,
		IMPLIED_BY_SYMB,

		INHIBITS_TEXT,
		INHIBITS_SYMB,
		INHIBITED_BY_TEXT,
		INHIBITED_BY_SYMB,

		LEFT_TEXT,
		LEFT_SYMB,
		RIGHT_TEXT,
		RIGHT_SYMB,

		NOT_LEFT_TEXT,
		NOT_LEFT_SYMB,
		NOT_RIGHT_TEXT,
		NOT_RIGHT_SYMB,
	)
	for _, test := range tests {
		res, err := FileParser.ParseString("", test.input)
		actual := res.Expressions[0].Bool.Rest.Op
		assert.NoError(t, err)
		assert.Equal(t, actual, test.expected)
	}
}
