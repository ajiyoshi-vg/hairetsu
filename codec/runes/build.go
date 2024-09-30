package runes

import (
	"io"
	"iter"

	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/codec/dict"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type FillableDict codec.FillableDict[rune]

func FromReadSeeker[D FillableDict](r io.ReadSeeker, f item.Factory, d D) error {
	return NewBuilder[D]().Build(r, f, d)
}

func FromSlice[D FillableDict](xs []string, f item.Factory, d D) error {
	return NewBuilder[D]().BuildSlice(xs, f, d)
}

func NewBuilder[D FillableDict]() *dict.Builder[rune, string, D] {
	return dict.NewBuilder(scan.Lines, runeSeq, newEncoder[D])
}

func newEncoder[D FillableDict](dict D) codec.Encoder[string] {
	return NewEncoder(dict)
}

func runeSeq(seq iter.Seq[string]) iter.Seq[rune] {
	return func(yield func(rune) bool) {
		for line := range seq {
			for _, c := range line {
				if !yield(c) {
					return
				}
			}
		}
	}
}
