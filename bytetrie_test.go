package hairetsu

import (
	"encoding/binary"
	"os"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/codec/bytes"
	"github.com/ajiyoshi-vg/hairetsu/codec/composer"
	"github.com/ajiyoshi-vg/hairetsu/codec/trie"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
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
		{
			title: "binary",
			ng: [][]byte{
				[]byte("hoge"),
			},
			prefix: []byte{1, 0, 0},
			num:    2,
			data:   generateBytes(65535),
		},
	}
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			// build index from data
			b := composer.NewBytes(bytes.NewIdentityDict())
			origin, err := b.ComposeFromSlice(c.data)
			assert.NoError(t, err)

			f, err := os.CreateTemp("", "bytes")
			assert.NoError(t, err)
			defer os.Remove(f.Name())

			// write index to file
			_, err = origin.WriteTo(f)
			assert.NoError(t, err)

			err = f.Close()
			assert.NoError(t, err)

			// open index from file
			file, err := trie.OpenFile(
				f.Name(),
				bytes.NewEncoder(bytes.NewIdentityDict()),
			)
			assert.NoError(t, err)

			// open index from mmap
			mmap, err := trie.OpenMmap(
				f.Name(),
				bytes.NewEncoder(bytes.NewIdentityDict()),
			)
			assert.NoError(t, err)
			defer func() { assert.NoError(t, mmap.Close()) }()

			kinds := map[string]trie.Searchable[[]byte]{
				"origin": origin.Searcher(),
				"file":   file.Searcher(),
				"mmap":   mmap.Searcher(),
			}

			for kind, da := range kinds {
				t.Run(kind, func(t *testing.T) {
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
				})
			}
		})
	}
}

func trimLeft(b []byte) []byte {
	if len(b) < 2 {
		return b
	}
	if b[0] == 0 {
		return trimLeft(b[1:])
	}
	return b
}

func generateBytes(num uint32) [][]byte {
	ret := make([][]byte, 0, num)
	for i := 0; i < int(num); i++ {
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(i))
		buf = trimLeft(buf)
		ret = append(ret, buf)
	}
	return ret
}
