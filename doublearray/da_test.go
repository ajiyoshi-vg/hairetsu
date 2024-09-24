package doublearray

import (
	"bytes"
	"os"
	"slices"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestDoubleArrayReadWrite(t *testing.T) {
	cases := []struct {
		title  string
		data   []item.Item
		ng     []word.Word
		prefix word.Word
		num    int
	}{
		{
			title: "keytree",
			data: item.FromWords(
				word.Word{5, 4, 3},
				word.Word{5, 4, 3, 2, 1},
			),
			ng: []word.Word{
				{5},
				{5, 4},
				{5, 4, 3, 2},
			},
			prefix: word.Word{5, 4, 3, 2, 1},
			num:    2,
		},
	}

	for _, c := range cases {
		origin, err := NewBuilder().StreamBuild(slices.Values(c.data))
		assert.NoError(t, err)

		tmp, err := os.CreateTemp("", "test")
		assert.NoError(t, err)
		defer os.Remove(tmp.Name())

		_, err = origin.WriteTo(tmp)
		assert.NoError(t, err)
		assert.NoError(t, tmp.Close())

		tmp, err = os.Open(tmp.Name())
		assert.NoError(t, err)

		restored1 := New()
		_, err = restored1.ReadFrom(tmp)
		assert.NoError(t, err)
		assert.NoError(t, tmp.Close())

		copied := FromArray(origin.Array())

		mmaped, err := OpenMmap(tmp.Name())
		assert.NoError(t, err)

		restored2 := New()
		buf := &bytes.Buffer{}
		_, err = mmaped.WriteTo(buf)
		assert.NoError(t, err)
		_, err = restored2.ReadFrom(buf)
		assert.NoError(t, err)

		type trie interface {
			ExactMatchSearch(word.Word) (node.Index, error)
			CommonPrefixSearch(word.Word) ([]node.Index, error)
			Nodes
		}

		das := []trie{origin, restored1, copied, mmaped, restored2}

		for _, da := range das {
			s := GetStat(da)
			assert.Equal(t, len(c.data), s.Leaf)

			for _, x := range c.data {
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
		}
	}
}
