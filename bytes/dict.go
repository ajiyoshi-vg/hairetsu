package bytes

import (
	"encoding"
	"fmt"
	"io"
	"math"
	"sort"

	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Dict []byte

var (
	_ encoding.BinaryMarshaler   = (Dict)(nil)
	_ encoding.BinaryUnmarshaler = (*Dict)(nil)
)

type Builder struct {
	byteCount map[byte]uint32
}

func New(bs []byte) Dict {
	b := NewBuilder()
	b.Add(bs)
	return b.Build()
}

func (d Dict) Code(b byte) word.Code {
	return word.Code(d[b])
}

func (d Dict) Word(bs []byte) word.Word {
	ret := make(word.Word, 0, len(bs))
	for _, b := range bs {
		ret = append(ret, d.Code(b))
	}
	return ret
}

func (d Dict) WithNameSpace(ns, key []byte) word.Word {
	ret := make(word.Word, 0, len(ns)+len(key)+1)
	for _, b := range ns {
		ret = append(ret, d.Code(b))
	}
	ret = append(ret, word.Separator)
	for _, b := range key {
		ret = append(ret, d.Code(b))
	}
	return ret
}

func (d Dict) MarshalBinary() ([]byte, error) {
	ret := make([]byte, len(d))
	copy(ret, d)
	return ret, nil
}

func (d *Dict) UnmarshalBinary(bs []byte) error {
	if len(bs) != math.MaxUint8 {
		return fmt.Errorf("want %d bytes got %d", math.MaxUint8, len(bs))
	}
	*d = make([]byte, math.MaxUint8)
	copy(*d, bs)
	return nil
}

func (d Dict) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(d[:])
	return int64(n), err
}

func (d *Dict) ReadFrom(r io.Reader) (int64, error) {
	*d = make([]byte, math.MaxUint8)
	n, err := r.Read(*d)
	return int64(n), err
}

func NewBuilder() *Builder {
	return &Builder{
		byteCount: map[byte]uint32{},
	}
}

func (x *Builder) Add(bs []byte) {
	for _, b := range bs {
		x.byteCount[b] += 1
	}
}

func (x *Builder) Build() Dict {
	type tmp struct {
		b byte
		n uint32
	}

	buf := make([]tmp, 0, len(x.byteCount))
	for b, n := range x.byteCount {
		buf = append(buf, tmp{b: b, n: n})
	}

	sort.Slice(buf, func(i, j int) bool {
		return buf[i].n > buf[j].n
	})

	var ret [math.MaxUint8]byte
	for i, x := range buf {
		ret[x.b] = byte(i)
	}
	return Dict(ret[:])
}

func FromReader(r io.Reader) (Dict, error) {
	b := NewBuilder()
	buf := make([]byte, math.MaxUint8)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			b.Add(buf[:n])
		}
		if err == io.EOF {
			return b.Build(), nil
		}
		if err != nil {
			return nil, err
		}
	}
}
