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

func FromReadSeeker[T FillableDict](r io.ReadSeeker, f item.Factory, d T) error {
	b := dict.NewBuilder(scan.Lines, runeSeq, newEncoder[T])
	return b.Build(r, f, d)
}

func newEncoder[T FillableDict](dict T) codec.Encoder[string] {
	return NewEncoder(dict)
}

func runeSeq(r io.Reader) iter.Seq[rune] {
	return func(yield func(rune) bool) {
		for line := range scan.Lines(r) {
			for _, c := range line {
				if !yield(c) {
					return
				}
			}
		}
	}
}
