package hairetsu

import (
	"testing"

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
			data := c.data
			da, err := newRuneTrieBuilder().Build(data)
			assert.NoError(t, err)

			for i, x := range data {
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

			s := da.data.Stat()
			filled := float64(s.Length-s.NumEmpty) / float64(s.Length)
			t.Logf("size:%d filled:%f", s.Length, filled)
		})
	}
}
