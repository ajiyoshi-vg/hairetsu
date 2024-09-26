package doublebyte

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Dict interface {
	codec.Dict[uint16, word.Code]
	io.WriterTo
}
type inverseDict codec.Dict[word.Code, uint16]

var (
	Identity Dict = (*identity)(nil)
	_        Dict = mapDict{}
	_        Dict = ArrayDict{}
)

type identity struct{}
type inverseIdentity struct{}

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

type mapDict map[uint16]word.Code
type inverseMapDict map[word.Code]uint16

func MapDict(m map[uint16]word.Code) mapDict {
	return m
}

func (m mapDict) Code(x uint16) word.Code {
	ret, ok := m[x]
	if !ok {
		return word.Unknown
	}
	return ret
}
func (m mapDict) Inverse() codec.Dict[word.Code, uint16] {
	ret := make(inverseMapDict, len(m))
	for k, v := range m {
		ret[v] = k
	}
	return ret
}
func (m mapDict) WriteTo(w io.Writer) (int64, error) {
	a := NewArrayDict(m)
	return a.WriteTo(w)
}
func (m inverseMapDict) Code(x word.Code) uint16 {
	return m[x]
}
func (m inverseMapDict) Inverse() codec.Dict[uint16, word.Code] {
	ret := make(mapDict, len(m))
	for k, v := range m {
		ret[v] = k
	}
	return ret
}

type ArrayDict [math.MaxUint16]word.Code

func NewArrayDict(m mapDict) ArrayDict {
	ret := [math.MaxUint16]word.Code{}
	for n := range math.MaxUint16 {
		ret[n] = m.Code(uint16(n))
	}
	return ret
}

func (a ArrayDict) Code(x uint16) word.Code {
	return a[x]
}
func (a ArrayDict) Inverse() codec.Dict[word.Code, uint16] {
	ret := make(inverseMapDict, len(a))
	for n, c := range a {
		ret[c] = uint16(n)
	}
	return ret
}
func (a ArrayDict) WriteTo(w io.Writer) (int64, error) {
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
	if err != nil {
		return int64(n), err
	}
	bf := bytes.NewReader(buf)
	for i := 0; i < math.MaxUint16; i++ {
		var c uint32
		if err := binary.Read(bf, binary.LittleEndian, &c); err != nil {
			return int64(n), err
		}
		a[i] = word.Code(c)
	}
	return int64(n), nil
}
