package hairetsu

import (
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/codec/composer"
	"github.com/ajiyoshi-vg/hairetsu/codec/runes"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/stretchr/testify/assert"
)

func TestRuneTrieSearch(t *testing.T) {
	cases := []struct {
		title  string
		data   []string
		ng     []string
		prefix string
		num    int
	}{
		{
			title: "インドネシア",
			ng: []string{
				("hoge"),
				("印"),
			},
			prefix: ("印度尼西亚啊"),
			num:    2,
			data: []string{
				("印度"),
				("印度尼西亚"),
				("印加帝国"),
				("瑞士"),
				("瑞典"),
				("巴基斯坦"),
				("巴勒斯坦"),
				("以色列"),
				("巴比伦"),
				("土耳其"),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			encs := map[string]composer.Composable[string]{
				"rune-id":  composer.NewRunes(runes.NewIdentityDict()),
				"rune-map": composer.NewRunes(runes.NewMapDict()),
			}
			for name, enc := range encs {
				t.Run(name, func(t *testing.T) {
					x, err := enc.ComposeFromSlice(c.data)
					assert.NoError(t, err)

					da := x.Searcher()

					for i, x := range c.data {
						actual, err := da.ExactMatchSearch(x)
						assert.NoError(t, err, x)
						assert.Equal(t, node.Index(i), actual)
					}
					for _, x := range c.ng {
						_, err := da.ExactMatchSearch(x)
						assert.Error(t, err, x)
					}

					is, err := da.CommonPrefixSearch(c.prefix)
					assert.NoError(t, err)
					assert.Equal(t, c.num, len(is))

					t.Log(x.Stat())
				})
			}
		})
	}
}
