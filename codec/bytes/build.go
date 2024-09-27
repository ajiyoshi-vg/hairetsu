package bytes

import (
	"io"
	"iter"
	"slices"

	"github.com/ajiyoshi-vg/external/emit"
	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/codec/dict"
)

type FillableDict codec.FillableDict[byte]

func FromReadSeeker[T FillableDict](r io.ReadSeeker, f dict.Factory, d T) error {
	b := dict.NewBuilder(scan.ByteLines, byteSeq, newEncoder[T])
	return b.Build(r, f, d)
}

func newEncoder[T FillableDict](dict T) codec.Encoder[[]byte] {
	return NewEncoder[T](dict)
}

func byteSeq(r io.Reader) iter.Seq[byte] {
	return func(yield func(byte) bool) {
		for line := range scan.ByteLines(r) {
			emit.All(slices.Values(line), yield)
		}
	}
}
