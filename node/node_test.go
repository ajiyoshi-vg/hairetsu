package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHoge(t *testing.T) {
	x := New(1)
	assert.Equal(t, Index(0), x.GetOffset())
	x.SetOffset(42)
	assert.Equal(t, Index(42), x.GetOffset())

	assert.False(t, x.IsTerminal())
	x.Terminate()
	assert.True(t, x.IsTerminal())

	assert.False(t, x.HasParent())
	x.SetParent(42)
	assert.True(t, x.HasParent())
	assert.True(t, x.IsChildOf(42))
}
