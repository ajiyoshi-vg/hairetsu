package node

import (
	"testing"
	"testing/quick"

	"github.com/ajiyoshi-vg/hairetsu/word"
)

func TestIndex(t *testing.T) {
	mustInverse := func(i Index, c word.Code) bool {
		actual := i.Forward(c).Backward(c)
		return i == actual
	}

	c := &quick.Config{
		MaxCountScale: 10000,
	}
	if err := quick.Check(mustInverse, c); err != nil {
		t.Error(err)
	}
}

func TestOneToOne(t *testing.T) {
	oneToOne := func(i Index, a, b word.Code) bool {
		if a == b {
			return true
		}
		x, y := i.Forward(a), i.Forward(b)
		return x != y
	}

	c := &quick.Config{
		MaxCountScale: 10000,
	}
	if err := quick.Check(oneToOne, c); err != nil {
		t.Error(err)
	}
}

func TestOneToOne2(t *testing.T) {
	oneToOne := func(a, b Index, c word.Code) bool {
		if a == b {
			return true
		}
		x, y := a.Forward(c), b.Forward(c)
		return x != y
	}

	c := &quick.Config{
		MaxCountScale: 10000,
	}
	if err := quick.Check(oneToOne, c); err != nil {
		t.Error(err)
	}
}
