package doublearray

import (
	"bytes"
	"slices"
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/runes"
	"github.com/ajiyoshi-vg/hairetsu/token"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestDoubleArraySearch(t *testing.T) {
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

		das := []Nodes{origin}

		for _, da := range das {
			s := GetStat(da)
			assert.Equal(t, len(c.data), s.Leaf)

			for _, x := range c.data {
				actual, err := Words{}.ExactMatchSearch(da, x.Word)
				assert.NoError(t, err)
				assert.Equal(t, node.Index(x.Val), actual)

				bs, err := x.Word.Bytes()
				assert.NoError(t, err)
				actual, err = Bytes{}.ExactMatchSearch(da, bs)
				assert.NoError(t, err)
				assert.Equal(t, node.Index(x.Val), actual)
			}

			for _, x := range c.ng {
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

func TestRunesDictSearch(t *testing.T) {
	cases := []struct {
		title  string
		data   string
		ng     []string
		prefix string
		num    int
	}{
		{
			title: "rune dict",
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
		r := bytes.NewBufferString(c.data)
		ks, dict, err := runes.FromWalker(token.NewLinedString(r))
		assert.NoError(t, err)

		err = NewBuilder().Build(origin, ks)
		assert.NoError(t, err)

		das := []Nodes{origin}

		for _, da := range das {
			assert.Equal(t, ks.LeafNum(), GetStat(da).Leaf)

			ss := strings.Split(c.data, "\n")
			for _, s := range ss {
				_, err := RunesDict(dict).ExactMatchSearch(da, s)
				assert.NoError(t, err)
			}

			for _, x := range c.ng {
				_, err := RunesDict(dict).ExactMatchSearch(da, x)
				assert.Error(t, err)
			}

			actual, err := RunesDict(dict).CommonPrefixSearch(da, c.prefix)
			assert.NoError(t, err)
			assert.Equal(t, c.num, len(actual))
		}
	}
}

func TestBytesDictSearch(t *testing.T) {
	cases := []struct {
		title  string
		data   string
		ng     []string
		prefix string
		num    int
	}{
		{
			title: "bytes dict",
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
		r := bytes.NewBufferString(c.data)
		ks, dict, err := runes.FromWalker(token.NewLinedString(r))
		assert.NoError(t, err)

		err = NewBuilder().Build(origin, ks)
		assert.NoError(t, err)

		das := []Nodes{origin}

		for _, da := range das {
			assert.Equal(t, ks.LeafNum(), GetStat(da).Leaf)

			ss := strings.Split(c.data, "\n")
			for _, s := range ss {
				_, err := RunesDict(dict).ExactMatchSearch(da, s)
				assert.NoError(t, err)
			}

			for _, x := range c.ng {
				_, err := RunesDict(dict).ExactMatchSearch(da, x)
				assert.Error(t, err)
			}

			actual, err := RunesDict(dict).CommonPrefixSearch(da, c.prefix)
			assert.NoError(t, err)
			assert.Equal(t, c.num, len(actual))
		}
	}
}
