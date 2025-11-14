package dicescript

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalDiceBasic(t *testing.T) {
	res, err := EvalDice("1+1")
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.Equal(t, "1+1", res.Expression)
		assert.Equal(t, "", res.Rest)
		assert.True(t, valueEqual(res.Value, ni(2)))
		assert.Equal(t, "2", res.ValueText)
	}
}

func TestEvalDiceWithOptions(t *testing.T) {
	res, err := EvalDice("E5", func(ctx *Context) error {
		return ctx.RegCustomDice(`E(\d+)`, func(ctx *Context, groups []string, _ any) (*VMValue, string, error) {
			return ni(42), "custom:E", nil
		})
	})
	assert.NoError(t, err)
	if assert.NotNil(t, res) {
		assert.True(t, valueEqual(res.Value, ni(42)))
		assert.Contains(t, res.Detail, "custom:E")
	}
}

func TestEvalDiceInvalidInput(t *testing.T) {
	_, err := EvalDice("")
	assert.Error(t, err)
}
