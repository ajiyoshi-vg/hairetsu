package doublearray

import (
	"slices"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestStreamBuild(t *testing.T) {
	cases := map[string]struct {
		input  []item.Item
		ng     []word.Word
		prefix word.Word
		num    int
		option []Option
	}{
		"normal": {
			input: item.FromWords(
				word.Word{5, 4, 3},
				word.Word{5, 4, 3, 2, 1},
			),
			ng: []word.Word{
				{5},
				{5, 4},
				{5, 4, 3, 2},
			},
			prefix: word.Word{5, 4, 3, 2, 1, 2, 3, 4, 5},
			num:    2,
		},
		"no val": {
			input: []item.Item{
				{Word: word.Word{5, 4, 3}},
				{Word: word.Word{5, 4, 3, 2, 1}},
			},
			ng: []word.Word{
				{5},
				{5, 4},
				{5, 4, 3, 2},
			},
			prefix: word.Word{5, 4, 3, 2, 1, 2, 3, 4, 5},
			num:    2,
		},
		"chunked": {
			input: item.FromWords(
				word.Word{5, 4, 3},
				word.Word{5, 4, 3, 2, 1},
			),
			ng: []word.Word{
				{5},
				{5, 4},
				{5, 4, 3, 2},
			},
			prefix: word.Word{5, 4, 3, 2, 1, 2, 3, 4, 5},
			num:    2,
			option: []Option{
				StreamChunkSize(1),
			},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			b := NewBuilder(c.option...)
			da, err := b.StreamBuild(slices.Values(c.input))
			assert.NoError(t, err)

			s := GetStat(da)
			assert.Equal(t, len(c.input), s.Leaf)

			for _, x := range c.input {
				actual, err := da.ExactMatchSearch(x.Word)
				assert.NoError(t, err)
				assert.Equal(t, node.Index(x.Val), actual)
			}

			for _, x := range c.ng {
				_, err := da.ExactMatchSearch(x)
				assert.Error(t, err)
			}

			actual, err := da.CommonPrefixSearch(c.prefix)
			assert.NoError(t, err)
			assert.Equal(t, c.num, len(actual))
		})
	}
}
