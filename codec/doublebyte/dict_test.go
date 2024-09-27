package doublebyte

import (
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestArrayDict(t *testing.T) {
	cases := map[string]struct {
		input   []byte
		unknown []byte
	}{
		"normal": {
			input:   []byte{1, 2, 3},
			unknown: []byte{4, 5, 6},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			dict := instantBuild(NewArrayDict(), c.input)
			enc := NewEncoder(dict)

			for code := range enc.Iter(c.input) {
				assert.NotEqual(t, word.Unknown, code)
			}

			for code := range enc.Iter(c.unknown) {
				if code == word.Unknown {
					continue
				}
				assert.Equal(t, word.Backspace, code)
			}
		})
	}
}
