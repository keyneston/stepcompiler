package stepcompiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskIsState(t *testing.T) {
	assert.Implements(t, (*State)(nil), &Pass{})
}

func TestTask(t *testing.T) {
	step := NewBuilder()
	task := NewTask("Foo").Next(NewTask("Bar").Comment("Bar does bar things"))

	expected := `
{
    "StartAt": "Foo",
    "States": {
        "Bar": {
            "Type": "Task",
            "End": true,
			"Comment": "Bar does bar things"
        },
        "Foo": {
            "Type": "Task",
            "Next": "Bar"
        }
    }
}
`

	step.StartAt(task)

	output, err := step.Render()
	assert.NoError(t, err)
	assert.JSONEq(t, string(output), expected)
}
