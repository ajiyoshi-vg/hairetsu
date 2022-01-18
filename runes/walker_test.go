package runes

import (
	"bytes"
	"io"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/token"
	"github.com/stretchr/testify/assert"
)

func TestStringLines(t *testing.T) {
	cases := []struct {
		title  string
		input  io.Reader
		expect []string
	}{
		{
			title: "normal",
			input: bytes.NewBufferString("aaa\nbb\ncccc"),
			expect: []string{
				"aaa",
				"bb",
				"cccc",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			tree, dict, err := FromWalker(token.NewLinedString(c.input))
			assert.NoError(t, err)
			for _, s := range c.expect {
				w, err := dict.Word(s)
				assert.NoError(t, err)

				i, err := tree.Get(w)
				assert.NoError(t, err)
				assert.NotNil(t, i)
			}
		})
	}
}
