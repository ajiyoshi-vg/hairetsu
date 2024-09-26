package doublebyte

import (
	"io"

	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type Factory interface {
	Put(item.Item)
}

func FromReadSeeker[D Dict](dict D, r io.ReadSeeker, f Factory) (D, error) {
	b := NewBuilder(dict)
	for line := range scan.ByteLines(r) {
		b.Add(line)
	}
	b.Build()

	var zero D
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return zero, err
	}

	enc := NewEncoder(dict)
	var i uint32
	for line := range scan.ByteLines(r) {
		f.Put(item.New(enc.Encode(line), i))
		i++
	}
	return dict, nil
}

type Builder[D Dict] struct {
	counter map[uint16]int
	dest    D
}

func NewBuilder[D Dict](dest D) *Builder[D] {
	return &Builder[D]{
		counter: make(map[uint16]int),
		dest:    dest,
	}
}

func (b *Builder[D]) Add(x []byte) {
	for i := range DoubleBytes(x) {
		b.counter[i] += 1
	}
}

func (b *Builder[D]) Build() {
	b.dest.Fill(b.counter)
}
