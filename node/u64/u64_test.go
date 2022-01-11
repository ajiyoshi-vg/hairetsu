package u64

import (
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/stretchr/testify/assert"
)

func TestHoge(t *testing.T) {
	x := New(1)
	assert.Equal(t, node.Index(0), x.GetOffset())
	x.SetOffset(42)
	assert.Equal(t, node.Index(42), x.GetOffset())

	assert.False(t, x.IsTerminal())
	x.Terminate()
	assert.True(t, x.IsTerminal())

	x.SetParent(42)
	assert.Equal(t, node.Index(42), x.GetParent())
	assert.True(t, x.HasParent())
	assert.True(t, x.IsChildOf(42))
}
