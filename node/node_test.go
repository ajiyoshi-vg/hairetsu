package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNode(t *testing.T) {
	var err error
	x := Root()

	// inital state
	assert.False(t, x.HasParent())
	assert.False(t, x.HasOffset())
	assert.False(t, x.IsTerminal())
	assert.False(t, x.IsUsed())
	_, err = x.GetNextEmptyNode()
	assert.NoError(t, err)
	_, err = x.GetPrevEmptyNode()
	assert.NoError(t, err)

	// can SetNext/SetPrev
	assert.NoError(t, x.SetNextEmptyNode(10))
	assert.NoError(t, x.SetPrevEmptyNode(10))

	// SetNext/SetPrev never changes state
	// still can GetNext/GetPrev
	_, err = x.GetNextEmptyNode()
	assert.NoError(t, err)
	_, err = x.GetPrevEmptyNode()
	assert.NoError(t, err)

	////////////////////////////////
	// SetOffset
	x.SetOffset(42)

	// SetOffset changes state
	assert.False(t, x.HasParent())
	assert.True(t, x.HasOffset()) // changed by SetOffset
	assert.False(t, x.IsTerminal())
	assert.True(t, x.IsUsed()) // changed by SetOffset
	_, err = x.GetNextEmptyNode()
	assert.NoError(t, err)
	_, err = x.GetPrevEmptyNode()
	assert.Error(t, err) // now can't GetPrev
	assert.NoError(t, x.SetNextEmptyNode(10))
	assert.Error(t, x.SetPrevEmptyNode(10)) // now can't SetPrev

	////////////////////////////////
	// SetParent
	x.SetParent(0)

	// SetParent changes state
	assert.True(t, x.HasParent()) // changed by SetParent
	assert.True(t, x.HasOffset())
	assert.False(t, x.IsTerminal())
	assert.True(t, x.IsUsed())
	_, err = x.GetNextEmptyNode()
	assert.Error(t, err) // now can't GetNext
	_, err = x.GetPrevEmptyNode()
	assert.Error(t, err)
	assert.Error(t, x.SetNextEmptyNode(10)) // now can't SetNext
	assert.Error(t, x.SetPrevEmptyNode(10))

	////////////////////////////////
	// Terminate
	x.Terminate()

	assert.True(t, x.HasParent())
	assert.True(t, x.HasOffset())
	assert.True(t, x.IsTerminal()) // changed by Terminate
	assert.True(t, x.IsUsed())
	_, err = x.GetNextEmptyNode()
	assert.Error(t, err)
	_, err = x.GetPrevEmptyNode()
	assert.Error(t, err)
	assert.Error(t, x.SetNextEmptyNode(10))
	assert.Error(t, x.SetPrevEmptyNode(10))
}
