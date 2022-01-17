package doublearray

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/keyset"
	"github.com/ajiyoshi-vg/hairetsu/keytree"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestDoubleArraySearch(t *testing.T) {
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

		das := []Nodes{origin}

		for _, da := range das {
			s := GetStat(da)
			assert.Equal(t, c.data.LeafNum(), s.Leaf)

			c.data.WalkLeaf(func(ws word.Word, val uint32) error {
				actual, err := ExactMatchSearchInterface(da, ws)
				assert.NoError(t, err)
				assert.Equal(t, node.Index(val), actual)

				actual, err = Words{}.ExactMatchSearch(da, ws)
				assert.NoError(t, err)
				assert.Equal(t, node.Index(val), actual)

				bs, err := ws.Bytes()
				assert.NoError(t, err)
				actual, err = Bytes{}.ExactMatchSearch(da, bs)
				assert.NoError(t, err)
				assert.Equal(t, node.Index(val), actual)

				return nil
			})

			for _, x := range c.ng {
				_, err = ExactMatchSearchInterface(da, x)
				assert.Error(t, err)

				_, err = Words{}.ExactMatchSearch(da, x)
				assert.Error(t, err)

				bs, err := x.Bytes()
				assert.NoError(t, err)
				_, err = Bytes{}.ExactMatchSearch(da, bs)
				assert.Error(t, err)
			}

			actual, err := Words{}.CommonPrefixSearch(da, c.prefix)
			assert.NoError(t, err)
			assert.Equal(t, c.num, len(actual))

			bs, err := c.prefix.Bytes()
			assert.NoError(t, err)
			actual, err = Bytes{}.CommonPrefixSearch(da, bs)
			assert.NoError(t, err)
			assert.Equal(t, c.num, len(actual))
		}
	}
}

func TestStringSearch(t *testing.T) {
	cases := []struct {
		title  string
		data   string
		ng     []string
		prefix string
		num    int
	}{
		{
			title: "keyset",
			data:  "aaa\nab\nabc",
			ng: []string{
				"abcc",
				"aa",
				"a",
			},
			prefix: "abcc",
			num:    2,
		},
	}

	for _, c := range cases {
		origin := New()
		ks, dict, err := keytree.FromStringLines(bytes.NewBufferString(c.data))
		assert.NoError(t, err)

		err = NewBuilder().Build(origin, ks)
		assert.NoError(t, err)

		das := []Nodes{origin}

		for _, da := range das {
			assert.Equal(t, ks.LeafNum(), GetStat(da).Leaf)

			ss := strings.Split(c.data, "\n")
			for _, s := range ss {
				_, err := Runes(dict).ExactMatchSearch(da, s)
				assert.NoError(t, err)
			}

			for _, x := range c.ng {
				_, err := Runes(dict).ExactMatchSearch(da, x)
				assert.Error(t, err)
			}

			actual, err := Runes(dict).CommonPrefixSearch(da, c.prefix)
			assert.NoError(t, err)
			assert.Equal(t, c.num, len(actual))
		}
	}
}
