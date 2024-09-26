package doublebyte

import (
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestDoubleByte(t *testing.T) {
	cases := map[string]struct {
		dict   codec.Dict[uint16, word.Code]
		input  []byte
		expect []byte
	}{
		"id:0": {
			dict:   Identity,
			input:  []byte{0},
			expect: []byte{0},
		},
		"id:0,1": {
			dict:   Identity,
			input:  []byte{0, 1},
			expect: []byte{0, 1},
		},
		"id:0,1,2": {
			dict:   Identity,
			input:  []byte{0, 1, 2},
			expect: []byte{0, 1, 2},
		},
		"map:1,2": {
			dict:   mapDict{0x0201: 1},
			input:  []byte{1, 2},
			expect: []byte{1, 2},
		},
		"map:1,2,3": {
			dict:   mapDict{0x0201: 1, 0x03: 2},
			input:  []byte{1, 2, 3},
			expect: []byte{1, 2, 3},
		},
		"array:1,2,3": {
			dict:   NewArrayDict(mapDict{0x0201: 1, 0x03: 2}),
			input:  []byte{1, 2, 3},
			expect: []byte{1, 2, 3},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			enc := NewEncoder(c.dict)
			dec := enc.Decoder()
			actual, err := dec.Decode(enc.Encode(c.input))
			assert.NoError(t, err)
			assert.Equal(t, c.input, actual)
		})
	}
}
