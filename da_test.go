package hairetsu

import (
	"bytes"
	"sort"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	cases := []struct {
		title  string
		data   [][]byte
		prefix []byte
		num    int
	}{
		{
			title:  "インドネシア",
			prefix: []byte("印度尼西亚啊"),
			num:    2,
			data: [][]byte{
				[]byte("印度"),
				[]byte("印度尼西亚"),
				[]byte("印加帝国"),
				[]byte("瑞士"),
				[]byte("瑞典"),
				[]byte("巴基斯坦"),
				[]byte("巴勒斯坦"),
				[]byte("以色列"),
				[]byte("巴比伦"),
				[]byte("土耳其"),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			data := c.data
			sort.Slice(data, func(i, j int) bool {
				return bytes.Compare(data[i], data[j]) < 0
			})
			builder := &builder{}
			da, err := builder.FromBytes(data)
			assert.NoError(t, err)

			for i, x := range data {
				actual, err := da.ExactMatchSearch(x)
				assert.NoError(t, err, x)
				assert.Equal(t, node.Index(i), actual)
			}

			is, err := da.CommonPrefixSearch(c.prefix)
			assert.NoError(t, err)
			assert.Equal(t, c.num, len(is))
		})
	}
}
