package runes

import (
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/codec/dict"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	cases := map[string]struct {
		input  string
		expect string
	}{
		"hello": {
			input:  "hello",
			expect: "hello",
		},
		"日本語": {
			input:  "日本語",
			expect: "日本語",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			kinds := map[string]WordDict{
				"map": NewMapDict(),
				"id":  NewIdentityDict(),
			}
			for kind, d := range kinds {
				t.Run(kind, func(t *testing.T) {
					r := strings.NewReader(c.input)
					dict := dict.InstantCount(d, runeSeq(r))
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
