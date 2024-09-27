package dict

import (
	"bytes"
	"encoding/binary"
	"io"
	"sort"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"golang.org/x/exp/constraints"
)

type Map[T constraints.Integer] map[T]word.Code
type inverseMap[T constraints.Integer] map[word.Code]T

var (
	_ codec.WordDict[int]        = (Map[int])(nil)
	_ codec.Dict[word.Code, int] = (inverseMap[int])(nil)
)

func (m Map[T]) Code(x T) word.Code {
	ret, ok := m[x]
	if !ok {
		return word.Unknown
	}
	return ret
}

func (m Map[T]) Fill(count map[T]int) {
	m.fill(count)
}

func (m Map[T]) fill(count map[T]int) Map[T] {
	type touple struct {
		i T
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

func (m Map[T]) Inverse() codec.Dict[word.Code, T] {
	ret := make(inverseMap[T], len(m))
	for k, v := range m {
		ret[v] = k
	}
	return ret
}

func (m Map[T]) WriteTo(w io.Writer) (int64, error) {
	buf := &bytes.Buffer{}
	size := len(m) * 8
	if err := binary.Write(buf, binary.LittleEndian, uint32(size)); err != nil {
		return 0, err
	}
	for x, c := range m {
		if err := binary.Write(buf, binary.LittleEndian, uint32(x)); err != nil {
			return 0, err
		}
		if err := binary.Write(buf, binary.LittleEndian, uint32(c)); err != nil {
			return 0, err
		}
	}
	return io.Copy(w, buf)
}

func (m Map[T]) ReadFrom(r io.Reader) (int64, error) {
	var ret int64
	sizeBuf := make([]byte, 4)
	n, err := io.ReadFull(r, sizeBuf)
	ret += int64(n)
	if err != nil {
		return ret, err
	}

	var size uint32
	if _, err := binary.Decode(sizeBuf, binary.LittleEndian, &size); err != nil {
		return ret, err
	}

	buf := make([]byte, size)
	n, err = io.ReadFull(r, buf)
	ret += int64(n)
	if err != nil {
		return ret, err
	}
	br := bytes.NewReader(buf)
	for {
		var x, c uint32
		if err := binary.Read(br, binary.LittleEndian, &x); err != nil {
			if err == io.EOF {
				return ret, nil
			}
			return ret, err
		}
		if err := binary.Read(br, binary.LittleEndian, &c); err != nil {
			if err == io.EOF {
				return ret, nil
			}
			return ret, err
		}
		m[T(x)] = word.Code(c)
	}
}

func (m inverseMap[T]) Code(x word.Code) T {
	ret, ok := m[x]
	if !ok {
		var zero T
		return zero
	}
	return ret
}

func (m inverseMap[T]) Inverse() codec.Dict[T, word.Code] {
	ret := make(Map[T], len(m))
	for k, v := range m {
		ret[v] = k
	}
	return ret
}
