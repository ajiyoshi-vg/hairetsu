package bytes

import (
	"bytes"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	cases := []struct {
		title  string
		input  []byte
		expect map[byte]word.Code
	}{
		{
			title: "normal",
			input: []byte("aabbbc"),
			expect: map[byte]word.Code{
				'b': 0,
				'a': 1,
				'c': 2,
			},
		},
	}

	for _, c := range cases {
		original := New(c.input)
		for b, c := range c.expect {
			assert.Equal(t, c, original.Code(b))
		}

		t.Run("FromReader", func(t *testing.T) {
			restored, err := FromReader(bytes.NewBuffer(c.input))
			assert.NoError(t, err)
			for b, c := range c.expect {
				assert.Equal(t, c, restored.Code(b))
			}
		})

		t.Run("MarshalText/UnmarshalText", func(t *testing.T) {
			tmp, err := original.MarshalBinary()
			assert.NoError(t, err)

			restored := Dict{}
			assert.NoError(t, restored.UnmarshalBinary(tmp))
			assert.Equal(t, original, restored)
		})

		t.Run("WriteTo/ReadFrom", func(t *testing.T) {
			buf := &bytes.Buffer{}
			n, err := original.WriteTo(buf)
			assert.NoError(t, err)
			assert.Equal(t, int64(buf.Len()), n)

			restored := Dict{}
			m, err := restored.ReadFrom(bytes.NewReader(buf.Bytes()))
			assert.NoError(t, err)
			assert.Equal(t, m, n)

			assert.Equal(t, original, restored)
		})
	}
}

func TestWord(t *testing.T) {
	cases := []struct {
		title  string
		input  []byte
		expect word.Word
	}{
		{
			title:  "normal",
			input:  []byte("abbbac"),
			expect: word.Word{1, 0, 0, 0, 1, 2},
		},
	}

	for _, c := range cases {
		dict := New(c.input)
		assert.Equal(t, c.expect, dict.Word(c.input))
	}
}
