package boolean

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBooleanParsing(t *testing.T) {
	tests := []struct {
		input   string
		literal *string
		ident   *string
	}{
		{"true", ptr("true"), nil},
		{"false", ptr("false"), nil},
		{"A", nil, ptr("A")},
		{"X1", nil, ptr("X1")},
	}
	for _, test := range tests {
		expr, err := Parser.ParseString("", test.input)
		assert.NoError(t, err)

		if test.literal != nil {
			assert.NotNil(t, expr.Literal)
			assert.Equal(t, *test.literal, expr.Literal.Value)
		} else {
			assert.Nil(t, expr.Literal)
		}

		if test.ident != nil {
			assert.NotNil(t, expr.Ident)
			assert.Equal(t, *test.ident, expr.Ident.Name)
		} else {
			assert.Nil(t, expr.Ident)
		}
	}
}

func ptr(s string) *string {
	return &s
}
