package doublearray

import (
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
		t.Run(c.title, func(t *testing.T) {
			da := New()
			ks := keytree.FromWord(c.content)
			assert.NoError(t, NewBuilder().Build(da, ks))

			num := 0
			err := ForEach(da, func(actual word.Word, val uint32) error {
				assert.Equal(t, c.content[val], actual)
				num++
				return nil
			})
			assert.NoError(t, err)
			assert.Equal(t, len(c.content), num)
		})
	}
}
