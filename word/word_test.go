package word

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromByte(t *testing.T) {
	cases := []struct {
		title  string
		source []byte
		expect []byte
	}{
		{
			title:  "simple",
			source: []byte("日本語"),
			expect: []byte("日本語"),
		},
		{
			title:  "nil",
			source: nil,
			expect: []byte{},
		},
		{
			title:  "[]byte{0, math.MaxUint8}",
			source: []byte{0, math.MaxUint8},
			expect: []byte{0, math.MaxUint8},
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			actual, err := FromBytes(c.source).Bytes()
			assert.NoError(t, err)
			assert.Equal(t, c.expect, actual)
		})
	}
}

func TestNameSpace(t *testing.T) {
	cases := []struct {
		title string
		lhs   Word
		rhs   Word
		equal bool
	}{
		{
			title: "ns/key == ns/key",
			lhs:   WithNameSpace([]byte("ns"), []byte("key")),
			rhs:   WithNameSpace([]byte("ns"), []byte("key")),
			equal: true,
		},
		{
			title: "ns/key != nsk/ey",
			lhs:   WithNameSpace([]byte("ns"), []byte("key")),
			rhs:   WithNameSpace([]byte("nsk"), []byte("ey")),
			equal: false,
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			if c.equal {
				assert.Equal(t, c.lhs, c.rhs)
			} else {
				assert.NotEqual(t, c.lhs, c.rhs)
			}
		})
	}
}

func TestAt(t *testing.T) {
	cases := []struct {
		title  string
		input  Word
		expect []Code
	}{
		{
			title:  "EOS",
			input:  Word{0, 1, 2, math.MaxUint8, 4},
			expect: []Code{0, 1, 2, math.MaxUint8, 4, EOS, EOS},
		},
		{
			title:  "SEP",
			input:  WithNameSpace([]byte{0, 1, 2}, []byte{3, 4}),
			expect: []Code{0, 1, 2, Separator, 3, 4, EOS},
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			for i, expect := range c.expect {
				assert.Equal(t, expect, c.input.At(i))
			}
		})
	}
}
