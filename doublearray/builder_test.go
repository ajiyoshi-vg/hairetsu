package doublearray

import (
	"slices"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/keytree"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/schollz/progressbar"
	"github.com/stretchr/testify/assert"
)

func TestDoubleArray(t *testing.T) {
	cases := []struct {
		title  string
		data   walker
		ng     []word.Word
		prefix word.Word
		num    int
	}{
		{
			title: "keytree",
			data: fromWord([]word.Word{
				{5, 4, 3},
				{5, 4, 3, 2, 1},
			}),
			ng: []word.Word{
				{5},
				{5, 4},
				{5, 4, 3, 2},
			},
			prefix: word.Word{5, 4, 3, 2, 1, 2, 3, 4, 5},
			num:    2,
		},
	}
	for _, c := range cases {
		da := New()

		b := NewBuilder(OptionProgress(progressbar.New(0)))
		err := b.Build(da, c.data)
		assert.NoError(t, err)

		s := GetStat(da)
		assert.Equal(t, c.data.LeafNum(), s.Leaf)

		c.data.WalkLeaf(func(key word.Word, val uint32) error {
			actual, err := da.ExactMatchSearch(key)
			assert.NoError(t, err)
			assert.Equal(t, node.Index(val), actual)
			return nil
		})

		for _, x := range c.ng {
			_, err := da.ExactMatchSearch(x)
			assert.Error(t, err)
		}

		actual, err := da.CommonPrefixSearch(c.prefix)
		assert.NoError(t, err)
		assert.Equal(t, c.num, len(actual))
	}
}

func fromWord(data []word.Word) *keytree.Tree {
	ks := keytree.New()
	for i, x := range data {
		ks.Put(x, uint32(i))
	}
	return ks
}

func TestStreamBuild(t *testing.T) {
	cases := map[string]struct {
		input  []Item
		ng     []word.Word
		prefix word.Word
		num    int
	}{
		"normal": {
			input: []Item{
				{word.Word{5, 4, 3}, 1},
				{word.Word{5, 4, 3, 2, 1}, 2},
			},
			ng: []word.Word{
				{5},
				{5, 4},
				{5, 4, 3, 2},
			},
			prefix: word.Word{5, 4, 3, 2, 1, 2, 3, 4, 5},
			num:    2,
		},
		"no val": {
			input: []Item{
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
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			da := New()
			b := NewBuilder()
			err := b.StreamBuild(da, slices.Values(c.input))
			assert.NoError(t, err)

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
