package keytree

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func str(x *uint32) string {
	if x == nil {
		return "nil"
	}
	return fmt.Sprintf("%d", *x)
}

func TestWalkNode(t *testing.T) {
	words := []word.Word{
		{1, 2, 3},
		{1, 2, 3, 4, 5},
		{1, 2, 4},
	}
	x := FromWord(words)

	assert.Equal(t, len(words), x.LeafNum())
	for i, w := range words {
		actual, err := x.Get(w)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, uint32(i), *actual)
	}

	ss := []string{}
	err := x.WalkNode(func(pre word.Word, brs []word.Code, val *uint32) error {
		sort.Slice(brs, func(i, j int) bool {
			return brs[i] < brs[j]
		})
		ss = append(ss, fmt.Sprintf("(%v, %v, %s)", pre, brs, str(val)))
		return nil
	})
	assert.NoError(t, err)
	expect := `([1 2 3 4 5], [], 1)
([1 2 3 4], [5], nil)
([1 2 3], [4], 0)
([1 2 4], [], 2)
([1 2], [3 4], nil)
([1], [2], nil)
([], [1], nil)`

	sort.Strings(ss)
	actual := strings.Join(ss, "\n")
	assert.Equal(t, expect, actual)
}

func TestWalkLeaf(t *testing.T) {
	words := []word.Word{
		{1, 2, 3},
		{1, 2, 3, 4, 5},
		{1, 2, 4},
	}
	x := FromWord(words)
	assert.Equal(t, len(words), x.LeafNum())

	for i, w := range words {
		actual, err := x.Get(w)
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, uint32(i), *actual)
	}

	ss := []string{}
	err := x.WalkLeaf(func(pre word.Word, val uint32) error {
		ss = append(ss, fmt.Sprintf("(%v, %d)", pre, val))
		return nil
	})
	assert.NoError(t, err)
	expect := `([1 2 3 4 5], 1)
([1 2 3], 0)
([1 2 4], 2)`

	sort.Strings(ss)
	actual := strings.Join(ss, "\n")
	assert.Equal(t, expect, actual)
}
