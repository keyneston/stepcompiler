package step

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSucceedIsState(t *testing.T) {
	assert.Implements(t, (*State)(nil), &Succeed{})
}

func TestSucceed(t *testing.T) {
	step := NewBuilder()
	pass := NewPass("Foo").Next(NewSucceed("Success"))

	expected := `
{
    "StartAt": "Foo",
    "States": {
        "Success": {
            "Type": "Succeed"
        },
        "Foo": {
            "Type": "Pass",
            "Next": "Success"
        }
    }
}`
	step.StartAt(pass)

	output, err := step.Render()
	assert.NoError(t, err)
	assert.JSONEq(t, string(output), expected)
}
