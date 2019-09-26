package stepcompiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	step := NewBuilder()
	task := NewTask("Foo").Next(NewTask("Bar"))

	step.StartAt(task)

	output, err := step.Render()
	assert.NoError(t, err)
	assert.Equal(t, string(output), "THIS DOESN'T WORK YET")
}
