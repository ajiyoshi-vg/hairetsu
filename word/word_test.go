package word

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type checker func(i int) error

var positive = func(i int) error {
	if i > 0 {
		return nil
	}
	return fmt.Errorf("want positive got %d", i)
}
var negative = func(i int) error {
	if i < 0 {
		return nil
	}
	return fmt.Errorf("want negative got %d", i)
}
var equal = func(i int) error {
	if i == 0 {
		return nil
	}
	return fmt.Errorf("want zero got %d", i)
}

func TestCompare(t *testing.T) {
	cases := []struct {
		title   string
		lhs     Word
		rhs     Word
		checker checker
	}{
		{
			title:   "nil = nil",
			lhs:     nil,
			rhs:     nil,
			checker: equal,
		},
		{
			title:   "nil == {}",
			lhs:     nil,
			rhs:     Word{},
			checker: equal,
		},
		{
			title:   "{} = nil",
			lhs:     Word{},
			rhs:     nil,
			checker: equal,
		},
		{
			title:   "{} = {}",
			lhs:     Word{},
			rhs:     Word{},
			checker: equal,
		},
		{
			title:   "{1, 2} > {1, 1}",
			lhs:     Word{1, 2},
			rhs:     Word{1, 1},
			checker: positive,
		},
		{
			title:   "{1, 2, 3} > {1, 2}",
			lhs:     Word{1, 2, 3},
			rhs:     Word{1, 2},
			checker: positive,
		},
		{
			title:   "{1, 1} < {1, 2}",
			lhs:     Word{1, 1},
			rhs:     Word{1, 2},
			checker: negative,
		},
		{
			title:   "{1, 2} < {1, 2, 3}",
			lhs:     Word{1, 2},
			rhs:     Word{1, 2, 3},
			checker: negative,
		},
		{
			title:   "{1, 2} = {1, 2}",
			lhs:     Word{1, 2},
			rhs:     Word{1, 2},
			checker: equal,
		},
	}
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			actual := Compare(c.lhs, c.rhs)
			assert.NoError(t, c.checker(actual))
		})
	}
}

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
			actual, err := FromBytes(c.source).bytes()
			assert.NoError(t, err)
			assert.Equal(t, c.expect, actual)
		})
	}
}
