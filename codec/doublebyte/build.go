package doublebyte

import (
	"io"

	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type Factory interface {
	Put(item.Item)
}

func FromReadSeeker[D Dict](r io.ReadSeeker, f Factory, dict D) error {
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

type builder[D Dict] struct {
	count map[uint16]int
	dest  D
}

func newBuilder[D Dict](dest D) *builder[D] {
	return &builder[D]{
		count: make(map[uint16]int),
		dest:  dest,
	}
}

func (b *builder[D]) add(x []byte) {
	for i := range DoubleBytes(x) {
		b.count[i] += 1
	}
}

func (b *builder[D]) build() D {
	b.dest.Fill(b.count)
	return b.dest
}

func instantBuild[D Dict](dest D, data []byte) D {
	b := newBuilder(dest)
	b.add(data)
	return b.build()
}
