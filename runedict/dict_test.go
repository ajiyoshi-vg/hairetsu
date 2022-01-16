package runedict

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
		expect RuneDict
	}{
		{
			title: "normal",
			input: "aab\nbbc",
			expect: RuneDict{
				'b': 0,
				'a': 1,
				'c': 2,
			},
		},
	}

	for _, c := range cases {
		actual := New(strings.Split(c.input, "\n"))
		assert.Equal(t, c.expect, actual)

		restored0, err := FromLines(bytes.NewBufferString(c.input))
		assert.NoError(t, err)
		assert.Equal(t, c.expect, restored0)

		tmp, err := restored0.MarshalText()
		assert.NoError(t, err)

		restored1 := RuneDict{}
		assert.NoError(t, restored1.UnmarshalText(tmp))
		assert.Equal(t, c.expect, restored1)

		buf := &bytes.Buffer{}
		n, err := restored1.WriteTo(buf)
		assert.NoError(t, err)
		assert.Equal(t, int64(buf.Len()), n)

		restored2 := RuneDict{}
		m, err := restored2.ReadFrom(bytes.NewReader(buf.Bytes()))
		assert.NoError(t, err)
		assert.Equal(t, m, n)

		assert.Equal(t, c.expect, restored2)
	}
}

func TestCode(t *testing.T) {
	cases := []struct {
		title  string
		dict   RuneDict
		input  rune
		expect word.Code
	}{
		{
			title: "normal",
			dict: RuneDict{
				'a': 42,
			},
			input:  'a',
			expect: 42,
		},
		{
			title:  "unknown rune returns word.NONE",
			dict:   RuneDict{},
			input:  'a',
			expect: word.NONE,
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
		dict  RuneDict
		input string
		check func(word.Word, error)
	}{
		{
			title: "normal",
			dict: RuneDict{
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
			dict:  RuneDict{},
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
