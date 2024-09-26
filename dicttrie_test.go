package hairetsu

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/stretchr/testify/assert"
)

func TestDictTrieSearch(t *testing.T) {
	cases := []struct {
		title  string
		data   [][]byte
		ng     [][]byte
		prefix []byte
		num    int
	}{
		{
			title: "インドネシア",
			ng: [][]byte{
				[]byte("hoge"),
				[]byte("印"),
			},
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
			da, err := NewDictTrieBuilder().BuildFromSlice(c.data)
			assert.NoError(t, err)

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

			t.Log(doublearray.GetStat(da.data))
		})
	}
}

func TestDictTrieBuild(t *testing.T) {
	cases := []struct {
		title string
		data  string
		ng    [][]byte
	}{
		{
			title: "BuildLines",
			data:  "aaa\nbbb\nabc\nabb",
			ng: [][]byte{
				[]byte("hoge"),
				[]byte("印"),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			data := bytes.NewReader([]byte(c.data))
			origin, err := NewDictTrieBuilder().BuildFromLines(data)
			assert.NoError(t, err)

			buf := &bytes.Buffer{}
			_, err = origin.WriteTo(buf)
			assert.NoError(t, err)

			restored := &DictTrie{}
			_, err = restored.ReadFrom(buf)
			assert.NoError(t, err)
			assert.Equal(t,
				doublearray.GetStat(origin.data),
				doublearray.GetStat(restored.data),
			)

			das := []*DictTrie{origin, restored}

			for _, da := range das {
				for i, x := range strings.Split(c.data, "\n") {
					actual, err := da.ExactMatchSearch([]byte(x))
					assert.NoError(t, err, x)
					assert.Equal(t, node.Index(i), actual)
				}
				for _, x := range c.ng {
					id, err := da.ExactMatchSearch(x)
					assert.Error(t, err, string(x), id)
				}
				t.Log(doublearray.GetStat(da.data))
			}
		})
	}
}
