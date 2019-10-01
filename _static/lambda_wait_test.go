package step

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLambdaWaitImplements(t *testing.T) {
	assert.Implements(t, (*State)(nil), &LambdaWait{})
	assert.Implements(t, (*ChainableState)(nil), &LambdaWait{})
}
