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

func FromReadSeeker[D FillableDict](r io.ReadSeeker, f item.Factory, d D) error {
	return NewBuilder[D]().Build(r, f, d)
}

func FromSlice[D FillableDict](xs [][]byte, f item.Factory, d D) error {
	return NewBuilder[D]().BuildSlice(xs, f, d)
}

func NewBuilder[D FillableDict]() *dict.Builder[byte, []byte, D, *Encoder[D]] {
	return dict.NewBuilder(scan.ByteLines, byteSeq, NewEncoder[D])
}

func NewBuilderWith[D FillableDict](items func(io.Reader) iter.Seq[[]byte]) *dict.Builder[byte, []byte, D, *Encoder[D]] {
	return dict.NewBuilder(items, byteSeq, NewEncoder[D])
}

func byteSeq(seq iter.Seq[[]byte]) iter.Seq[byte] {
	return func(yield func(byte) bool) {
		for line := range seq {
			emit.All(slices.Values(line), yield)
		}
	}
}
