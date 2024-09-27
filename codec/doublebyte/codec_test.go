package doublebyte

import (
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/stretchr/testify/assert"
)

func TestDoubleByte(t *testing.T) {
	cases := map[string]struct {
		dict   WordDict
		input  []byte
		expect []byte
	}{
		"id:0": {
			dict:   &codec.Identity[uint16]{},
			input:  []byte{0},
			expect: []byte{0},
		},
		"id:0,1": {
			dict:   &codec.Identity[uint16]{},
			input:  []byte{0, 1},
			expect: []byte{0, 1},
		},
		"id:0,1,2": {
			dict:   &codec.Identity[uint16]{},
			input:  []byte{0, 1, 2},
			expect: []byte{0, 1, 2},
		},
		"map:1,2": {
			dict:   codec.MapDict[uint16]{},
			input:  []byte{1, 2},
			expect: []byte{1, 2},
		},
		"map:1,2,3": {
			dict:   codec.MapDict[uint16]{},
			input:  []byte{1, 2, 3},
			expect: []byte{1, 2, 3},
		},
		"array:1,2": {
			dict:   codec.NewArrayDict[uint16](),
			input:  []byte{1, 2},
			expect: []byte{1, 2},
		},
		"array:1,2,3": {
			dict:   codec.NewArrayDict[uint16](),
			input:  []byte{1, 2, 3},
			expect: []byte{1, 2, 3},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			dict := codec.InstantCount(c.dict, DoubleBytes(c.input))
			enc := NewEncoder(dict)
			dec := enc.Decoder()
			actual, err := dec.Decode(enc.Encode(c.input))
			assert.NoError(t, err)
			assert.Equal(t, c.input, actual)
		})
	}
}
