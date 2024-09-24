package doublearray

import (
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestForeach(t *testing.T) {
	cases := []struct {
		title   string
		content []word.Word
	}{
		{
			content: []word.Word{
				{1, 2, 3, 4},
				{0, 1, 2},
				{0, 1, 2, 3, 4, 5},
			},
		},
		{
			content: []word.Word{
				word.FromBytes([]byte{1, 2, 3}),
				word.WithNameSpace([]byte{1, 2}, []byte{3}),
			},
		},
	}

	for _, c := range cases {
		origin, err := StreamBuild(slices.Values(item.FromWordSlice(c.content)))
		assert.NoError(t, err)

		tmp, err := os.CreateTemp("", "test")
		assert.NoError(t, err)
		defer os.Remove(tmp.Name())

		_, err = origin.WriteTo(tmp)
		assert.NoError(t, err)
		assert.NoError(t, tmp.Close())

		mmaped, err := OpenMmap(tmp.Name())
		assert.NoError(t, err)

		das := map[string]Nodes{
			"DoubleArray": origin,
			"Mmap":        mmaped,
		}

		for name, da := range das {
			t.Run(fmt.Sprintf("%s:%s", c.title, name), func(t *testing.T) {
				num := 0
				err := WalkLeaf(da, func(actual word.Word, val uint32) error {
					assert.Equal(t, c.content[val], actual)
					num++
					return nil
				})
				assert.NoError(t, err)
				assert.Equal(t, len(c.content), num)
			})
		}
	}
}
