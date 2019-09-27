package stepcompiler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDynamoGetImplements(t *testing.T) {
	assert.Implements(t, (*State)(nil), &DynamoGet{})
	assert.Implements(t, (*ChainableState)(nil), &DynamoGet{})
}

func TestDynamoGet(t *testing.T) {
	step := NewBuilder().StartAt(NewDynamoGet("DynamoGet").Timeout(time.Minute * 5))

	expected := `
{
    "StartAt": "DynamoGet",
    "States": {
        "DynamoGet": {
            "Type": "Task",
			"TimeoutSeconds": 300,
			"Resource": "arn:aws:states:::dynamodb:getItem",
			"End": true
        }
	}
}
`

	output, err := step.Render()
	assert.NoError(t, err)
	assert.JSONEq(t, string(output), expected)
}
