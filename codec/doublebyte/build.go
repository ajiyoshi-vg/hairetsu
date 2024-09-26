package doublebyte

import (
	"io"
	"sort"

	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Factory interface {
	Put(item.Item)
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

	enc := NewEncoder(dict)
	var i uint32
	for line := range scan.ByteLines(r) {
		f.Put(item.New(enc.Encode(line), i))
		i++
	}
	return dict, nil
}

type Builder struct {
	counter map[uint16]int
}

func NewBuilder() *Builder {
	return &Builder{
		counter: make(map[uint16]int),
	}
}

func (b *Builder) Add(x []byte) {
	for i := range DoubleBytes(x) {
		b.counter[i] += 1
	}
}

func (b *Builder) Build() mapDict {
	type touple struct {
		i uint16
		n int
	}

	buf := make([]touple, 0, len(b.counter))
	for i, n := range b.counter {
		buf = append(buf, touple{i: i, n: n})
	}

	sort.Slice(buf, func(i, j int) bool {
		return buf[i].n > buf[j].n
	})

	ret := make(mapDict, len(buf))
	for i, x := range buf {
		ret[x.i] = word.Code(i)
	}

	return ret
}
