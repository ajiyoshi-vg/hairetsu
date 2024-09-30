package bytes

import (
	"slices"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/codec/dict"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
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
				"map":   NewMapDict(),
				"array": NewArrayDict(),
				"id":    NewIdentityDict(),
			}
			for kind, d := range kinds {
				t.Run(kind, func(t *testing.T) {
					r := slices.Values([][]byte{c.input})
					dict := dict.InstantCount(d, byteSeq(r))
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
