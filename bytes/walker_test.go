package bytes

import (
	"fmt"
	"math"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

type mock struct {
	data []item.Item
}

func (m *mock) Put(x item.Item) {
	m.data = append(m.data, x)
}

func (m *mock) Get(w word.Word) error {
	for _, x := range m.data {
		if word.Compare(x.Word, w) == 0 {
			return nil
		}
	}
	return fmt.Errorf("not found")
}

func TestFromSlice(t *testing.T) {
	cases := []struct {
		title string
		input [][]byte
		num   int
	}{
		{
			title: "normal",
			input: [][]byte{
				{0, 1, 2},
				{0, math.MaxUint8, 4},
				{5, 7, 3},
				{5, 7, 3, 1},
			},
			num: 4,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			mock := &mock{}
			dict, err := FromSlice(c.input, mock)
			assert.NoError(t, err)

			i := 0
			for _, x := range c.input {
				assert.NoError(t, mock.Get(dict.Word(x)))
				i++
			}
			assert.Equal(t, c.num, i)
		})
	}
}
