package stepcompiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDynamoDeleteImplements(t *testing.T) {
	assert.Implements(t, (*State)(nil), &DynamoDelete{})
	assert.Implements(t, (*ChainableState)(nil), &DynamoDelete{})
}
