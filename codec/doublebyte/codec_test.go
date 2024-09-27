package doublebyte

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoubleByte(t *testing.T) {
	cases := map[string]struct {
		dict   Dict
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
			dict:   MapDict{},
			input:  []byte{1, 2},
			expect: []byte{1, 2},
		},
		"map:1,2,3": {
			dict:   MapDict{},
			input:  []byte{1, 2, 3},
			expect: []byte{1, 2, 3},
		},
		"array:1,2": {
			dict:   NewArrayDict(),
			input:  []byte{1, 2},
			expect: []byte{1, 2},
		},
		"array:1,2,3": {
			dict:   NewArrayDict(),
			input:  []byte{1, 2, 3},
			expect: []byte{1, 2, 3},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			enc := NewEncoder(instantBuild(c.dict, c.input))
			dec := enc.Decoder()
			actual, err := dec.Decode(enc.Encode(c.input))
			assert.NoError(t, err)
			assert.Equal(t, c.input, actual)
		})
	}
}
