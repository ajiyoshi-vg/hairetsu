package bytes

import (
	"io"
	"iter"
	"slices"

	"github.com/ajiyoshi-vg/external/emit"
	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/codec/dict"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type FillableDict codec.FillableDict[byte]

var Builder = dict.NewBuilder(scan.ByteLines, byteSeq, newEncoder[FillableDict])

func NewBuilder[D FillableDict]() *dict.Builder[byte, []byte, D] {
	return dict.NewBuilder(scan.ByteLines, byteSeq, newEncoder[D])
}

func FromReadSeeker[D FillableDict](r io.ReadSeeker, f item.Factory, d D) error {
	return Builder.Build(r, f, d)
}

func FromSlice[D FillableDict](xs [][]byte, f item.Factory, d D) error {
	b := dict.NewBuilder(scan.ByteLines, byteSeq, newEncoder[D])
	return b.BuildSlice(xs, f, d)
}

func newEncoder[D FillableDict](dict D) codec.Encoder[[]byte] {
	return NewEncoder(dict)
}

func byteSeq(seq iter.Seq[[]byte]) iter.Seq[byte] {
	return func(yield func(byte) bool) {
		for line := range seq {
			emit.All(slices.Values(line), yield)
		}
	}
}
