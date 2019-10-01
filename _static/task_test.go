package step

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskIsState(t *testing.T) {
	assert.Implements(t, (*State)(nil), &Task{})
	assert.Implements(t, (*ChainableState)(nil), &Task{})
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

// TestTaskCatchGathered tests that a function that is only included as a Next
// in a catch clause gets gathered and rendered into the entire state.
func TestTaskCatchGathered(t *testing.T) {
	step := NewBuilder().StartAt(
		NewTask("Foo").Catch(
			NewCatchClause().Next(
				NewTask("HandleError").Next(NewTask("HandleError2")),
			),
		),
	)

	expected := `
{
    "StartAt": "Foo",
    "States": {
        "Foo": {
            "Type": "Task",
			"End": true,
			"Catch": [{
				"ErrorEquals": ["States.ALL"],
				"Next": "HandleError"
			}]
        },
		"HandleError": {
			"Type": "Task",
			"Next": "HandleError2"
		},
		"HandleError2": {
			"Type": "Task",
			"End": true
		}
    }
}
`

	output, err := step.Render()
	assert.NoError(t, err)
	assert.JSONEq(t, string(output), expected)
}

func TestTaskTimeout(t *testing.T) {
	step := NewBuilder().StartAt(NewTask("Foo").Timeout(time.Minute * 5))

	expected := `
{
    "StartAt": "Foo",
    "States": {
        "Foo": {
            "Type": "Task",
			"TimeoutSeconds": 300,
			"End": true
        }
	}
}
`

	output, err := step.Render()
	assert.NoError(t, err)
	assert.JSONEq(t, expected, string(output))
}
