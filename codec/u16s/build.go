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

func FromReadSeeker[T FillableDict](r io.ReadSeeker, f item.Factory, d T) error {
	b := dict.NewBuilder(scan.ByteLines, uint16Seq, newEncoder[T])
	return b.Build(r, f, d)
}

func newEncoder[T FillableDict](dict T) codec.Encoder[[]byte] {
	return NewEncoder(dict)
}

func uint16Seq(r io.Reader) iter.Seq[uint16] {
	return func(yield func(uint16) bool) {
		for line := range scan.ByteLines(r) {
			emit.All(DoubleBytes(line), yield)
		}
	}
}
