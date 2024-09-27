package doublebyte

import (
	"io"

	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type Factory interface {
	Put(item.Item)
}

func FromReadSeeker[T FillableDict](r io.ReadSeeker, f Factory, dict T) error {
	b := codec.NewCounter(dict)
	for line := range scan.ByteLines(r) {
		b.Add(DoubleBytes(line))
	}
	b.Build()

	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return err
	}

	enc := NewEncoder(dict)
	var i uint32
	for line := range scan.ByteLines(r) {
		f.Put(item.New(enc.Encode(line), i))
		i++
	}
	return nil
}
