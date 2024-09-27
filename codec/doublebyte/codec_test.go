package doublebyte

import (
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/stretchr/testify/assert"
)

func TestDoubleByte(t *testing.T) {
	cases := map[string]struct {
		input  []byte
		expect []byte
	}{
		"0": {
			input:  []byte{0},
			expect: []byte{0},
		},
		"0,1": {
			input:  []byte{0, 1},
			expect: []byte{0, 1},
		},
		"0,1,2": {
			input:  []byte{0, 1, 2},
			expect: []byte{0, 1, 2},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			kinds := map[string]WordDict{
				"map":   codec.MapDict[uint16]{},
				"array": codec.NewArrayDict[uint16](),
				"id":    &codec.Identity[uint16]{},
			}
			for kind, d := range kinds {
				t.Run(kind, func(t *testing.T) {
					dict := codec.InstantCount(d, DoubleBytes(c.input))
					enc := NewEncoder(dict)
					dec := enc.Decoder()
					actual, err := dec.Decode(enc.Encode(c.input))
					assert.NoError(t, err)
					assert.Equal(t, c.expect, actual)
				})
			}
		})
	}
}
