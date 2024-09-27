package dict

import (
	"bufio"
	"bytes"
	"io"
	"slices"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestReaderWriter(t *testing.T) {
	cases := map[string]struct {
		input []uint16
	}{
		"normal": {
			input: []uint16{1, 2, 3},
		},
		"empty": {
			input: nil,
		},
	}

	for name, c := range cases {
		kinds := map[string]codec.WordDict[uint16]{
			"map":   Map[uint16]{},
			"array": NewArray[uint16](),
			"id":    NewIdentity[uint16](),
		}
		t.Run(name, func(t *testing.T) {
			for kind, d := range kinds {
				t.Run(kind, func(t *testing.T) {
					dict := InstantCount(d, slices.Values(c.input))

					for _, b := range c.input {
						code := dict.Code(b)
						assert.NotEqual(t, word.Unknown, code)
					}

					buf := &bytes.Buffer{}
					nWrite, err := dict.WriteTo(buf)
					assert.NoError(t, err)

					_, err = buf.WriteString("dummy body")
					assert.NoError(t, err)

					br := bufio.NewReader(buf)

					nRead, err := dict.ReadFrom(br)
					assert.NoError(t, err)
					assert.Equal(t, nWrite, nRead)

					body, err := io.ReadAll(br)
					assert.NoError(t, err)
					assert.Equal(t, "dummy body", string(body))

					for _, b := range c.input {
						code := dict.Code(b)
						assert.NotEqual(t, word.Unknown, code)
					}
				})
			}
		})
	}
}
