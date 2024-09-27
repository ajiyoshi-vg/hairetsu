package codec

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"unsafe"

	"github.com/ajiyoshi-vg/hairetsu/word"
	"golang.org/x/exp/constraints"
)

type tinyInteger interface {
	~int8 | ~int16 | ~uint8 | ~uint16
}

type ArrayDict[T constraints.Integer] []word.Code

var (
	_ WordDict[byte] = (ArrayDict[byte])(nil)
)

func NewArrayDict[T constraints.Integer](opt ...Option[ArrayDict[T]]) ArrayDict[T] {
	var ret ArrayDict[T]
	ret = make([]word.Code, ret.bufferLength())
	for _, f := range opt {
		f(&ret)
	}
	return ret
}

func (a ArrayDict[T]) Code(x T) word.Code {
	return a[x]
}

func (a ArrayDict[T]) Fill(count map[T]int) {
	a.fill(count)
}

func (a ArrayDict[T]) fill(count map[T]int) ArrayDict[T] {
	tmp := MapDict[T]{}
	tmp.fill(count)
	for n := range a.bufferLength() {
		a[n] = tmp.Code(T(n))
	}
	return a
}

func (a ArrayDict[T]) Inverse() Dict[word.Code, T] {
	ret := make(inverseMapDict[T], len(a))
	for i, c := range a {
		if c != word.Unknown {
			ret[c] = T(i)
		}
	}
	return ret
}

func (a ArrayDict[T]) WriteTo(w io.Writer) (int64, error) {
	buf := &bytes.Buffer{}
	for _, c := range a {
		err := binary.Write(buf, binary.LittleEndian, uint32(c))
		if err != nil {
			return 0, err
		}
	}
	return io.Copy(w, buf)
}

func (a ArrayDict[T]) ReadFrom(r io.Reader) (int64, error) {
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

func (a ArrayDict[T]) bufferLength() int {
	size := unsafe.Sizeof(T(0))
	switch size {
	case 1:
		return math.MaxUint8
	case 2:
		return math.MaxUint16
	default:
		panic("unsupported type")
	}
}
func wordSize() int {
	return int(unsafe.Sizeof(word.Code(0)))
}
