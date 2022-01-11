package doublearray

import (
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/keyset"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestDoubleArray(t *testing.T) {
	da := &DoubleArray{
		nodes: make([]node.Node, 10),
	}

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
	da := &DoubleArray{
		nodes: make([]node.Node, 5),
	}
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
	da := &DoubleArray{
		nodes: make([]node.Node, 5),
	}

	data := []word.Word{
		word.Word{1},
		word.Word{1, 2},
		word.Word{2, 3, 4, 5},
	}
	err := NewBuilder().Build(da, keyset.New(data))
	assert.NoError(t, err)

	s := da.Stat()
	assert.Equal(t, len(data), s.Leaf)

	expect := []string{
		"{base:0, next:10}",  // 0
		"{base:3, check:0}#", // 1
		"{base:3, check:0}",  // 2
		"{base:0, check:1}",  // 3
		"{base:1, check:5}",  // 4
		"{base:4, check:1}#", // 5
		"{base:3, check:2}",  // 6
		"{base:3, check:6}",  // 7
		"{base:9, check:7}#", // 8
		"{base:2, check:8}",  // 9
		"{prev:0, next:11}",  // 10
	}

	for i, want := range expect {
		assert.Equal(t, want, da.nodes[i].String(), i)
	}

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