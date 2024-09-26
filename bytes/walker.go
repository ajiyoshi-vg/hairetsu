package bytes

import (
	"io"

	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type Factory interface {
	Put(item.Item)
}

func FromSlice(xs [][]byte, f Factory) (Dict, error) {
	b := NewBuilder()
	for _, x := range xs {
		b.Add(x)
	}
	dict := b.Build()

	for i, x := range xs {
		w := dict.Word(x)
		f.Put(item.New(w, uint32(i)))
	}
	return dict, nil
}

func FromReadSeeker(r io.ReadSeeker, f Factory) (Dict, error) {
	b := NewBuilder()
	for line := range scan.ByteLines(r) {
		b.Add(line)
	}
	dict := b.Build()

	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	var i uint32
	for line := range scan.ByteLines(r) {
		w := dict.Word(line)
		f.Put(item.New(w, i))
		i++
	}
	return dict, nil
}
