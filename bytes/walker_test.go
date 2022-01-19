package bytes

import (
	"bytes"
	"math"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/token"
	"github.com/stretchr/testify/assert"
)

func TestWalker(t *testing.T) {
	cases := []struct {
		title string
		input string
		num   int
	}{
		{
			title: "normal",
			input: "aaa\nbcc\naba",
			num:   3,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			w := token.NewLinedBytes(bytes.NewBufferString(c.input))
			ks, dict, err := FromWalker(w)
			assert.NoError(t, err)
			assert.NoError(t, checkAllCodeDiffers(dict))
			i := 0
			err = w.Walk(func(x []byte) error {
				_, err := ks.Get(dict.Word(x))
				assert.NoError(t, err)
				i++
				return nil
			})
			assert.NoError(t, err)
			assert.Equal(t, c.num, i)
		})
	}
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
				[]byte{0, 1, 2},
				[]byte{0, math.MaxUint8, 4},
				[]byte{5, 7, 3},
				[]byte{5, 7, 3, 1},
			},
			num: 4,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			ks, dict, err := FromSlice(c.input)
			assert.NoError(t, err)

			i := 0
			for _, x := range c.input {
				_, err := ks.Get(dict.Word(x))
				assert.NoError(t, err)
				i++
			}
			assert.Equal(t, c.num, i)
		})
	}
}
