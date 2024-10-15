package u16s

import (
	"io"
	"iter"

	"github.com/ajiyoshi-vg/external/emit"
	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/codec/dict"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type FillableDict codec.FillableDict[uint16]

func FromReadSeeker[D FillableDict](r io.ReadSeeker, f item.Factory, d D) error {
	return NewBuilder[D]().Build(r, f, d)
}

func FromSlice[D FillableDict](xs [][]byte, f item.Factory, d D) error {
	return NewBuilder[D]().BuildSlice(xs, f, d)
}

func NewBuilder[D FillableDict]() *dict.Builder[uint16, []byte, D, *Encoder[D]] {
	return dict.NewBuilder(scan.ByteLines, uint16Seq, NewEncoder[D])
}
func NewBuilderWith[D FillableDict](
	items func(io.Reader) iter.Seq[[]byte],
) *dict.Builder[uint16, []byte, D, *Encoder[D]] {
	return dict.NewBuilder(items, uint16Seq, NewEncoder[D])
}

func uint16Seq(seq iter.Seq[[]byte]) iter.Seq[uint16] {
	return func(yield func(uint16) bool) {
		for line := range seq {
			emit.All(DoubleBytes(line), yield)
		}
	}
}
