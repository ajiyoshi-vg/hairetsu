package doublebyte

import (
	"io"

	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type Factory interface {
	Put(item.Item)
}

func FromReadSeeker[T FillableDict](r io.ReadSeeker, f Factory, dict T) error {
	b := newBuilder(dict)
	for line := range scan.ByteLines(r) {
		b.add(line)
	}
	b.build()

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

type builder[T Fillable] struct {
	count map[uint16]int
	dest  T
}

func newBuilder[T Fillable](dest T) *builder[T] {
	return &builder[T]{
		count: make(map[uint16]int),
		dest:  dest,
	}
}

func (b *builder[T]) add(x []byte) {
	for i := range DoubleBytes(x) {
		b.count[i] += 1
	}
}

func (b *builder[T]) build() T {
	b.dest.Fill(b.count)
	return b.dest
}

func instantBuild[T Fillable](dest T, data []byte) T {
	b := newBuilder(dest)
	b.add(data)
	return b.build()
}
