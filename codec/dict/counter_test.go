package dict

import (
	"slices"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
	cases := map[string]struct {
		input   []uint16
		unknown []uint16
	}{
		"normal": {
			input:   []uint16{1, 2, 3},
			unknown: []uint16{4, 5, 6},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			kinds := map[string]codec.FillableDict[uint16]{
				"map":   Map[uint16]{},
				"array": NewArray[uint16](),
			}
			for kind, d := range kinds {
				t.Run(kind, func(t *testing.T) {
					dict := InstantCount(d, slices.Values(c.input))

					for _, b := range c.input {
						code := dict.Code(b)
						assert.NotEqual(t, word.Unknown, code)
					}

					for _, b := range c.unknown {
						code := dict.Code(b)
						assert.Equal(t, word.Unknown, code)
					}
				})
			}
			t.Run("identity", func(t *testing.T) {
				dict := InstantCount(NewIdentity[uint16](), slices.Values(c.input))

				xs := append(c.input, c.unknown...)
				for _, b := range xs {
					code := dict.Code(b)
					assert.Equal(t, word.Code(b), code)
				}
			})
		})
	}
}
