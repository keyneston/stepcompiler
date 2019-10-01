package step

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassImplements(t *testing.T) {
	assert.Implements(t, (*State)(nil), &Pass{})
	assert.Implements(t, (*ChainableState)(nil), &Pass{})
}

func TestPass(t *testing.T) {
	step := NewBuilder()
	pass := NewPass("Foo").Next(NewPass("Bar"))

	expected := `
{
    "StartAt": "Foo",
    "States": {
        "Bar": {
            "Type": "Pass",
            "End": true
        },
        "Foo": {
            "Type": "Pass",
            "Next": "Bar"
        }
    }
}`
	step.StartAt(pass)

	output, err := step.Render()
	assert.NoError(t, err)
	assert.JSONEq(t, string(output), expected)
}
