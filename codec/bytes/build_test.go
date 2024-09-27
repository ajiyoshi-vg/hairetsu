package bytes

import (
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/codec/dict"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/stretchr/testify/assert"
)

type mock struct {
	item []item.Item
}

func (m *mock) Put(x item.Item) {
	m.item = append(m.item, x)
}

func TestBuild(t *testing.T) {
	cases := map[string]struct {
		input  string
		expect [][]byte
	}{
		"normal": {
			input: "hello\nworld\n",
			expect: [][]byte{
				[]byte("hello"),
				[]byte("world"),
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			kinds := map[string]WordDict{
				"map":   dict.MapDict[byte]{},
				"array": dict.NewArrayDict[byte](),
				"id":    &dict.Identity[byte]{},
			}
			for kind, d := range kinds {
				t.Run(kind, func(t *testing.T) {
					r := strings.NewReader(c.input)
					f := &mock{}

					assert.NoError(t, FromReadSeeker(r, f, d))

					dec := NewEncoder(d).Decoder()

					actual := make([][]byte, 0, len(f.item))
					for _, item := range f.item {
						x, err := dec.Decode(item.Word)
						assert.NoError(t, err)
						actual = append(actual, x)
					}
					assert.Equal(t, c.expect, actual)
				})
			}
		})
	}
}
