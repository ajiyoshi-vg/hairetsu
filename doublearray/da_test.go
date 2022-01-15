package doublearray

import (
	"bytes"
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

		buf := &bytes.Buffer{}
		_, err = origin.WriteTo(buf)
		assert.NoError(t, err)

		restored := New()
		_, err = restored.ReadFrom(buf)
		assert.NoError(t, err)

		copied := FromArray(origin.Array())

		das := []*DoubleArray{origin, restored, copied}

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
