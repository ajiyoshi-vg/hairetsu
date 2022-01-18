package runes

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	cases := []struct {
		title  string
		input  string
		expect Dict
	}{
		{
			title: "normal",
			input: "aab\nbbc",
			expect: Dict{
				'b': 0,
				'a': 1,
				'c': 2,
			},
		},
	}

	for _, c := range cases {
		actual := New(strings.Split(c.input, "\n"))
		assert.Equal(t, c.expect, actual)

		original, err := fromLines(bytes.NewBufferString(c.input))
		assert.NoError(t, err)
		assert.Equal(t, c.expect, original)

		t.Run("MarshalText/UnmarshalText", func(t *testing.T) {
			tmp, err := original.MarshalText()
			assert.NoError(t, err)

			restored := Dict{}
			assert.NoError(t, restored.UnmarshalText(tmp))
			assert.Equal(t, c.expect, restored)
		})

		t.Run("MarshalBinary/UnmarshalBinary", func(t *testing.T) {
			tmp, err := original.MarshalBinary()
			assert.NoError(t, err)

			restored := Dict{}
			assert.NoError(t, restored.UnmarshalBinary(tmp))
			assert.Equal(t, c.expect, restored)
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

			assert.Equal(t, c.expect, restored)
		})
	}
}

func TestCode(t *testing.T) {
	cases := []struct {
		title  string
		dict   Dict
		input  rune
		expect word.Code
	}{
		{
			title: "normal",
			dict: Dict{
				'a': 42,
			},
			input:  'a',
			expect: 42,
		},
		{
			title:  "unknown rune returns word.Unknown",
			dict:   Dict{},
			input:  'a',
			expect: word.Unknown,
		},
	}

	for _, c := range cases {
		actual := c.dict.Code(c.input)
		assert.Equal(t, c.expect, actual)
	}
}

func TestWord(t *testing.T) {
	cases := []struct {
		title string
		dict  Dict
		input string
		check func(word.Word, error)
	}{
		{
			title: "normal",
			dict: Dict{
				'a': 0,
				'b': 1,
				'c': 2,
			},
			input: "abbac",
			check: func(actual word.Word, err error) {
				assert.Equal(t, word.Word{0, 1, 1, 0, 2}, actual)
			},
		},
		{
			title: "unknwon code returns error",
			dict:  Dict{},
			input: "a",
			check: func(_ word.Word, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, c := range cases {
		c.check(c.dict.Word(c.input))
	}
}
