package doublearray

import (
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/keyset"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestDoubleArray(t *testing.T) {
	da := New(10)

	data := []word.Word{
		word.Word{5, 4, 3},
		word.Word{5, 4, 3, 2, 1},
	}
	err := NewBuilder().Build(da, keyset.New(data))
	assert.NoError(t, err)

	s := da.Stat()
	assert.Equal(t, len(data), s.Leaf)

	for i, x := range data {
		actual, err := da.ExactMatchSearch(x)
		assert.NoError(t, err)
		assert.Equal(t, node.Index(i), actual)
	}

	ng := []word.Word{
		word.Word{5},
		word.Word{5, 4},
		word.Word{5, 4, 3, 2},
	}
	for _, x := range ng {
		_, err := da.ExactMatchSearch(x)
		assert.Error(t, err)
	}
}

func TestInitDoubleArray(t *testing.T) {
	da := New(5)
	NewBuilder().init(da, 0)
	expect := []string{
		"{prev:0, next:1}",
		"{prev:0, next:2}",
		"{prev:1, next:3}",
		"{prev:2, next:4}",
		"{prev:3, next:5}",
	}
	for i, node := range da.nodes {
		assert.Equal(t, expect[i], node.String())
	}
	s := da.Stat()
	assert.Equal(t, 5, s.Size)
	assert.Equal(t, 0, s.Node)
	assert.Equal(t, 5, s.Empty)
}

func TestBuildDoubleArray(t *testing.T) {
	da := New(5)

	data := []word.Word{
		word.Word{1},
		word.Word{1, 2},
		word.Word{2, 3, 4, 5},
	}
	err := NewBuilder().Build(da, keyset.New(data))
	assert.NoError(t, err)

	s := da.Stat()
	assert.Equal(t, len(data), s.Leaf)

	for i, x := range data {
		actual, err := da.ExactMatchSearch(x)
		assert.NoError(t, err)
		assert.Equal(t, node.Index(i), actual)
	}

	ng := []word.Word{
		word.Word{1, 2, 3},
		word.Word{2},
		word.Word{2, 3},
		word.Word{2, 3, 4},
	}
	for _, x := range ng {
		_, err := da.ExactMatchSearch(x)
		assert.Error(t, err)
	}
}
