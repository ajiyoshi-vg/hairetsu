package item

import (
	"iter"

	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Factory interface {
	Put(Item)
}

type Item struct {
	Word word.Word
	Val  uint32
}

func New(word word.Word, val uint32) Item {
	return Item{Word: word, Val: val}
}

func FromWordSlice(xs []word.Word) []Item {
	ret := make([]Item, 0, len(xs))
	for i, x := range xs {
		ret = append(ret, New(x, uint32(i)))
	}
	return ret
}

func FromWords(xs ...word.Word) []Item {
	return FromWordSlice(xs)
}

func FromByteSlice(xs [][]byte) []Item {
	ret := make([]Item, 0, len(xs))
	for i, x := range xs {
		ret = append(ret, New(word.FromBytes(x), uint32(i)))
	}
	return ret
}

func FromBytes(xs ...[]byte) []Item {
	return FromByteSlice(xs)
}

func FromByteSeq(seq iter.Seq[[]byte]) iter.Seq[Item] {
	return func(yield func(Item) bool) {
		var i uint32
		for x := range seq {
			if !yield(New(word.FromBytes(x), i)) {
				return
			}
			i++
		}
	}
}
