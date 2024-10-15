package dict

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"unsafe"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"golang.org/x/exp/constraints"
)

type tinyInteger interface {
	~int8 | ~int16 | ~uint8 | ~uint16
}

type Array[T constraints.Integer] []word.Code

var (
	_ codec.WordDict[byte] = (Array[byte])(nil)
)

func NewArray[T constraints.Integer](opt ...Option[Array[T]]) Array[T] {
	var ret Array[T]
	ret = make([]word.Code, ret.bufferLength())
	for _, f := range opt {
		f(&ret)
	}
	return ret
}

func (a Array[T]) Code(x T) word.Code {
	return a[x]
}

func (a Array[T]) Fill(count map[T]int) {
	a.fill(count)
}

func (a Array[T]) fill(count map[T]int) Array[T] {
	tmp := Map[T]{}
	tmp.fill(count)
	for n := range a.bufferLength() {
		a[n] = tmp.Code(T(n))
	}
	return a
}

func (a Array[T]) Inverse() codec.Dict[word.Code, T] {
	ret := make(inverseMap[T], len(a))
	for i, c := range a {
		if c != word.Unknown {
			ret[c] = T(i)
		}
	}
	return ret
}

func (a Array[T]) WriteTo(w io.Writer) (int64, error) {
	buf := &bytes.Buffer{}
	for _, c := range a {
		err := binary.Write(buf, binary.LittleEndian, uint32(c))
		if err != nil {
			return 0, err
		}
	}
	return io.Copy(w, buf)
}

func (a Array[T]) ReadFrom(r io.Reader) (int64, error) {
	size := a.bufferLength()
	buf := make([]byte, size*wordSize())

	n, err := io.ReadFull(r, buf)
	ret := int64(n)
	if err != nil {
		return ret, err
	}
	bf := bytes.NewReader(buf)
	for i := 0; i < size; i++ {
		var c uint32
		if err := binary.Read(bf, binary.LittleEndian, &c); err != nil {
			return ret, err
		}
		a[i] = word.Code(c)
	}
	return ret, nil
}

func (a Array[T]) bufferLength() int {
	size := unsafe.Sizeof(T(0))
	switch size {
	case 1:
		return math.MaxUint8 + 1
	case 2:
		return math.MaxUint16 + 1
	default:
		panic("unsupported type")
	}
}
func wordSize() int {
	return int(unsafe.Sizeof(word.Code(0)))
}
