package keyset

import (
	"fmt"
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
	x := FromWord([]word.Word{
		word.Word{1, 2, 3},
		word.Word{1, 2, 3, 4, 5},
		word.Word{1, 2, 4},
	})

	ss := []string{"(prefix, branch)"}
	err := x.WalkNode(func(pre word.Word, brs []word.Code, val *uint32) error {
		ss = append(ss, fmt.Sprintf("(%v, %v, %s)", pre, brs, str(val)))
		return nil
	})
	assert.NoError(t, err)
	expect := `(prefix, branch)
([], [1 1 1], nil)
([1], [2 2 2], nil)
([1 2], [3 3 4], nil)
([1 2 3], [4], 0)
([1 2 3 4], [5], nil)
([1 2 3 4 5], [], 1)
([1 2 4], [], 2)`

	actual := strings.Join(ss, "\n")
	assert.Equal(t, expect, actual)
}

func TestWalkLeaf(t *testing.T) {
	x := FromWord([]word.Word{
		word.Word{1, 2, 3},
		word.Word{1, 2, 3, 4, 5},
		word.Word{1, 2, 4},
	})

	ss := []string{"(prefix, value)"}
	err := x.WalkLeaf(func(pre word.Word, val uint32) error {
		ss = append(ss, fmt.Sprintf("(%v, %d)", pre, val))
		return nil
	})
	assert.NoError(t, err)
	expect := `(prefix, value)
([1 2 3], 0)
([1 2 3 4 5], 1)
([1 2 4], 2)`

	actual := strings.Join(ss, "\n")
	assert.Equal(t, expect, actual)
}
