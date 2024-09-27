package doublebyte

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"sort"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Dict interface {
	codec.Dict[uint16, word.Code]
	Fill(map[uint16]int)
	io.WriterTo
	io.ReaderFrom
}
type inverseDict codec.Dict[word.Code, uint16]

var (
	Identity Dict = (*identity)(nil)
	_        Dict = MapDict{}
	_        Dict = (*ArrayDict)(nil)
)

type identity struct{}
type inverseIdentity struct{}

func (*identity) Fill(count map[uint16]int) {
}
func (*identity) Code(x uint16) word.Code {
	return word.Code(x)
}
func (*identity) Inverse() codec.Dict[word.Code, uint16] {
	return &inverseIdentity{}
}
func (*identity) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}
func (*identity) ReadFrom(r io.Reader) (int64, error) {
	return 0, nil
}
func (*inverseIdentity) Code(x word.Code) uint16 {
	return uint16(x)
}
func (*inverseIdentity) Inverse() codec.Dict[uint16, word.Code] {
	return &identity{}
}

type MapDict map[uint16]word.Code
type inverseMapDict map[word.Code]uint16

func (m MapDict) Fill(count map[uint16]int) {
	m.fill(count)
}
func (m MapDict) fill(count map[uint16]int) MapDict {
	type touple struct {
		i uint16
		n int
	}

	buf := make([]touple, 0, len(count))
	for i, n := range count {
		buf = append(buf, touple{i: i, n: n})
	}

	sort.Slice(buf, func(i, j int) bool {
		return buf[i].n > buf[j].n
	})

	for i, x := range buf {
		m[x.i] = word.Code(i)
	}
	return m
}
func (m MapDict) Code(x uint16) word.Code {
	ret, ok := m[x]
	if !ok {
		return word.Unknown
	}
	return ret
}
func (m MapDict) Inverse() codec.Dict[word.Code, uint16] {
	ret := make(inverseMapDict, len(m))
	for k, v := range m {
		ret[v] = k
	}
	return ret
}
func (m MapDict) WriteTo(w io.Writer) (int64, error) {
	return NewArrayDict(WithContent(m)).WriteTo(w)
}
func (m MapDict) ReadFrom(r io.Reader) (int64, error) {
	a := NewArrayDict()
	n, err := a.ReadFrom(r)
	if err != nil {
		return n, err
	}
	for i, c := range a {
		if c == word.Unknown {
			continue
		}
		m[uint16(i)] = c
	}
	return n, nil
}
func (m inverseMapDict) Code(x word.Code) uint16 {
	return m[x]
}
func (m inverseMapDict) Inverse() codec.Dict[uint16, word.Code] {
	ret := make(MapDict, len(m))
	for k, v := range m {
		ret[v] = k
	}
	return ret
}

type ArrayDict [math.MaxUint16]word.Code

func NewArrayDict(opt ...Option[ArrayDict]) *ArrayDict {
	ret := &ArrayDict{}
	for _, f := range opt {
		f(ret)
	}
	return ret
}

func (a *ArrayDict) Fill(count map[uint16]int) {
	a.fill(count)
}
func (a *ArrayDict) fill(count map[uint16]int) *ArrayDict {
	m := make(MapDict, len(count))
	m.Fill(count)
	*a = *NewArrayDict(WithContent(m))
	return a
}
func (a *ArrayDict) Code(x uint16) word.Code {
	return a[x]
}
func (a *ArrayDict) Inverse() codec.Dict[word.Code, uint16] {
	ret := make(inverseMapDict, len(a))
	for n, c := range a {
		ret[c] = uint16(n)
	}
	return ret
}
func (a *ArrayDict) WriteTo(w io.Writer) (int64, error) {
	buf := &bytes.Buffer{}
	for _, c := range a {
		err := binary.Write(buf, binary.LittleEndian, uint32(c))
		if err != nil {
			return 0, err
		}
	}
	return io.Copy(w, buf)
}
func (a *ArrayDict) ReadFrom(r io.Reader) (int64, error) {
	size := 4 * math.MaxUint16
	buf := make([]byte, size)
	n, err := io.ReadFull(r, buf)
	ret := int64(n)
	if err != nil {
		return ret, err
	}
	bf := bytes.NewReader(buf)
	for i := 0; i < math.MaxUint16; i++ {
		var c uint32
		if err := binary.Read(bf, binary.LittleEndian, &c); err != nil {
			return ret, err
		}
		a[i] = word.Code(c)
	}
	return ret, nil
}
