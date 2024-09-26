package hairetsu

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/doublearray"
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
			da, err := NewRuneTrieBuilder().BuildFromSlice(data)
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

			t.Log(doublearray.GetStat(da.data))
		})
	}
}

func TestRuneTrieBuild(t *testing.T) {
	cases := []struct {
		title string
		data  string
		ng    []string
	}{
		{
			title: "BuildLines",
			data:  "aaa\nbbb\nabc\nabb",
			ng: []string{
				("hoge"),
				("印"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			data := bytes.NewReader([]byte(c.data))
			origin, err := NewRuneTrieBuilder().BuildFromLines(data)
			assert.NoError(t, err)

			buf := &bytes.Buffer{}
			_, err = origin.WriteTo(buf)
			assert.NoError(t, err)

			restored := &RuneTrie{}
			_, err = restored.ReadFrom(buf)
			assert.NoError(t, err)

			das := []*RuneTrie{origin, restored}

			for _, da := range das {

				for i, x := range strings.Split(c.data, "\n") {
					actual, err := da.ExactMatchSearch(x)
					assert.NoError(t, err, x)
					assert.Equal(t, node.Index(i), actual)
				}
				for _, x := range c.ng {
					_, err := da.ExactMatchSearch(x)
					assert.Error(t, err, x)
				}
				t.Log(doublearray.GetStat(da.data))
			}
		})
	}
}
