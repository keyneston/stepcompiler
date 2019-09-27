package stepcompiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDynamoPutImplements(t *testing.T) {
	assert.Implements(t, (*State)(nil), &DynamoPut{})
	assert.Implements(t, (*ChainableState)(nil), &DynamoPut{})
}
