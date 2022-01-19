package doublearray

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/keytree"
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
				word.Word{1, 2, 3, 4},
				word.Word{0, 1, 2},
				word.Word{0, 1, 2, 3, 4, 5},
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
		origin := New()
		ks := keytree.FromWord(c.content)
		assert.NoError(t, NewBuilder().Build(origin, ks))

		tmp, err := ioutil.TempFile("", "test")
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
