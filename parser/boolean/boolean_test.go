package boolean

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrue(t *testing.T) {
	input := "true"
	expected := Ptr(TRUE)
	result, err := LXBoolParser.ParseString("", input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.LXBool.True)
}

func TestFalse(t *testing.T) {
	input := "false"
	expected := Ptr(FALSE)
	result, err := LXBoolParser.ParseString("", input)
	assert.NoError(t, err)
	assert.Equal(t, expected, result.LXBool.False)
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

func Ptr(s LXBoolString) *LXBoolString {
	return &s
}

func Deref(s *LXBoolString) LXBoolString {
	return *s
}
