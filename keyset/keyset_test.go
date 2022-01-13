package keyset

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestWalk(t *testing.T) {
	x := New([]word.Word{
		word.Word{1, 2, 3},
		word.Word{1, 2, 3, 4, 5},
		word.Word{1, 2, 4},
	})

	ss := []string{"(prefix, branches)"}
	err := x.WalkNode(func(pre word.Word, brs []word.Code, vals []uint32) error {
		ss = append(ss, fmt.Sprintf("(%v, %v)", pre, brs))
		return nil
	})
	assert.NoError(t, err)
	expect := `(prefix, branches)
([], [1 1 1])
([1], [2 2 2])
([1 2], [3 3 4])
([1 2 3], [0 4])
([1 2 3 4], [5])
([1 2 3 4 5], [0])
([1 2 4], [0])`

	actual := strings.Join(ss, "\n")
	assert.Equal(t, expect, actual)
}
