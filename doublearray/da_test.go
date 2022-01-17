package doublearray

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/keyset"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestDoubleArrayReadWrite(t *testing.T) {
	cases := []struct {
		title  string
		data   Walker
		ng     []word.Word
		prefix word.Word
		num    int
	}{
		{
			title: "keyset",
			data: keyset.FromWord([]word.Word{
				word.Word{5, 4, 3},
				word.Word{5, 4, 3, 2, 1},
			}),
			ng: []word.Word{
				word.Word{5},
				word.Word{5, 4},
				word.Word{5, 4, 3, 2},
			},
			prefix: word.Word{5, 4, 3, 2, 1},
			num:    2,
		},
	}

	for _, c := range cases {
		origin := New()

		err := NewBuilder().Build(origin, c.data)
		assert.NoError(t, err)

		tmp, err := ioutil.TempFile("", "test")
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
}
