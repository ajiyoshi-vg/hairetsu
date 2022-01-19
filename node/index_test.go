package node

import (
	"testing"
	"testing/quick"

	"github.com/ajiyoshi-vg/hairetsu/word"
)

func TestInverse(t *testing.T) {
	t.Run("always x.Forward(c).Backward(c) == x", func(t *testing.T) {
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
	})

	t.Run("always x.Forward(Special).Backward(Special) == x", func(t *testing.T) {
		mustInverse := func(i Index) bool {
			special := []word.Code{word.EOS, word.Separator}
			for _, s := range special {
				actual := i.Forward(s).Backward(s)
				if i != actual {
					return false
				}
			}
			return true
		}

		c := &quick.Config{
			MaxCountScale: 10000,
		}
		if err := quick.Check(mustInverse, c); err != nil {
			t.Error(err)
		}
	})
}

func TestOneToOne(t *testing.T) {
	t.Run("a != b <=> i.Forward(a) != i.Forward(b)", func(t *testing.T) {
		oneToOne := func(i Index, a, b word.Code) bool {
			if a == b {
				return i.Forward(a) == i.Forward(b)
			} else {
				return i.Forward(a) != i.Forward(b)
			}
		}

		c := &quick.Config{
			MaxCountScale: 10000,
		}
		if err := quick.Check(oneToOne, c); err != nil {
			t.Error(err)
		}
	})

	t.Run("a != b <=> a.Forward(c) != b.Forward(c)", func(t *testing.T) {
		oneToOne := func(a, b Index, c word.Code) bool {
			if a == b {
				return a.Forward(c) == b.Forward(c)
			} else {
				return a.Forward(c) != b.Forward(c)
			}
		}

		c := &quick.Config{
			MaxCountScale: 10000,
		}
		if err := quick.Check(oneToOne, c); err != nil {
			t.Error(err)
		}
	})

	t.Run("c != EOS <=> a.Forward(c) != a.Forward(EOS)", func(t *testing.T) {
		oneToOne := func(i Index, c word.Code) bool {
			if c == word.EOS {
				return i.Forward(c) == i.Forward(word.EOS)
			} else {
				return i.Forward(c) != i.Forward(word.EOS)
			}
		}

		c := &quick.Config{
			MaxCountScale: 10000,
		}
		if err := quick.Check(oneToOne, c); err != nil {
			t.Error(err)
		}
	})
}

func TestLabel(t *testing.T) {
	t.Run("child == offset.Forward(c) <=> GetLabel(offset, child) == c", func(t *testing.T) {
		oneToOne := func(offset, child Index, c word.Code) bool {
			if offset.Forward(c) == child {
				return Label(offset, child) == c
			} else {
				return Label(offset, child) != c
			}
		}

		c := &quick.Config{
			MaxCountScale: 10000,
		}
		if err := quick.Check(oneToOne, c); err != nil {
			t.Error(err)
		}
	})

	t.Run("child == offset.Forward(EOS) <=> GetLabel(offset, child) == EOS", func(t *testing.T) {
		oneToOne := func(offset, child Index) bool {
			if offset.Forward(word.EOS) == child {
				return Label(offset, child) == word.EOS
			} else {
				return Label(offset, child) != word.EOS
			}
		}

		c := &quick.Config{
			MaxCountScale: 10000,
		}
		if err := quick.Check(oneToOne, c); err != nil {
			t.Error(err)
		}
	})
}
