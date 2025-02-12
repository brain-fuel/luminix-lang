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
	input := " ( true ) "
	// expected := Ptr(TRUE)
	_, err := LXBoolParser.ParseString("", input)
	assert.NoError(t, err)
}

func Ptr(s LXBoolString) *LXBoolString {
	return &s
}

func Deref(s *LXBoolString) LXBoolString {
	return *s
}
