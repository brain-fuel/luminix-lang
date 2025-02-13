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
		_, err := LXBoolParser.ParseString("", test)
		errorStr := fmt.Sprintf("1:1: unexpected token \"%s\"", strings.Fields(test)[0])
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestTrue(t *testing.T) {
	input := "true"
	expected := Ptr(TRUE)
	result, err := LXBoolParser.ParseString("", input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.LXBool.True)
}

func TestTrueFailAfterTrueSucceed(t *testing.T) {
	tests := []string{
		"true t rue",
		"true tr ue",
		"true tru e",
		"true ture",
	}
	for _, test := range tests {
		_, err := LXBoolParser.ParseString("", test)
		errorStr := fmt.Sprintf("1:6: unexpected token \"%s\"", strings.Fields(test)[1])
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
		_, err := LXBoolParser.ParseString("", test)
		errorStr := fmt.Sprintf("2:1: unexpected token \"%s\"", strings.Fields(test)[1])
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
		_, err := LXBoolParser.ParseString("", test)
		errorStr := fmt.Sprintf("1:1: unexpected token \"%s\"", strings.Fields(test)[0])
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestFalse(t *testing.T) {
	input := "false"
	expected := Ptr(FALSE)
	result, err := LXBoolParser.ParseString("", input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.LXBool.False)
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
		_, err := LXBoolParser.ParseString("", test)
		errorStr := fmt.Sprintf("1:7: unexpected token \"%s\"", strings.Fields(test)[1])
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
		_, err := LXBoolParser.ParseString("", test)
		errorStr := fmt.Sprintf("2:1: unexpected token \"%s\"", strings.Fields(test)[1])
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
		_, err := LXBoolParser.ParseString("", test)
		errorStr := "1:7: unexpected token \"<EOF>\" (expected \")\")"
		expectedErr := errors.New(errorStr)
		assert.EqualError(t, err, expectedErr.Error())
	}
}

func TestParens(t *testing.T) {
	input := "(true)"
	expected := Ptr(TRUE)
	result, err := LXBoolParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.LXParBool)
	assert.Equal(t, expected, result.LXParBool.LXBoolExpr.LXBool.True)
}

func TestDoubleParens(t *testing.T) {
	input := "((true))"
	expected := Ptr(TRUE)
	result, err := LXBoolParser.ParseString("", input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.LXParBool)
	assert.NotNil(t, result.LXParBool.LXBoolExpr)
	assert.NotNil(t, result.LXParBool.LXBoolExpr.LXParBool)
	assert.NotNil(t, result.LXParBool.LXBoolExpr.LXParBool.LXBoolExpr)
	assert.NotNil(t, result.LXParBool.LXBoolExpr.LXParBool.LXBoolExpr.LXBool)
	assert.NotNil(t, result.LXParBool.LXBoolExpr.LXParBool.LXBoolExpr.LXBool.True)
	assert.Equal(t, expected, result.LXParBool.LXBoolExpr.LXParBool.LXBoolExpr.LXBool.True)
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
		actual, err := LXBoolParser.ParseString("", test)
		assert.NoError(t, err)
		assert.Equal(t, Ptr(TRUE), actual.LXParBool.LXBoolExpr.LXBool.True)
	}
}

func Ptr(s LXBoolString) *LXBoolString {
	return &s
}

func Deref(s *LXBoolString) LXBoolString {
	return *s
}
