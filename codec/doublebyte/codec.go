package doublebyte

import (
	"iter"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

var (
	_ codec.Encoder[[]byte] = (*Encoder[*ArrayDict])(nil)
	_ codec.Encoder[[]byte] = (*Encoder[MapDict])(nil)
	_ codec.Decoder[[]byte] = (*Decoder)(nil)
)

type Encoder[T Dict] struct {
	dictionary T
}

func NewEncoder[T Dict](d T) *Encoder[T] {
	return &Encoder[T]{
		dictionary: d,
	}
}

func DoubleBytes(x []byte) iter.Seq[uint16] {
	return func(yield func(uint16) bool) {
		for i := 0; i < len(x); i += 2 {
			val := uint16(x[i])
			if i+1 < len(x) {
				val |= uint16(x[i+1]) << 8
			}

			if !yield(val) {
				return
			}
		}
	}
}

func (enc Encoder[T]) Iter(x []byte) iter.Seq[word.Code] {
	return func(yield func(word.Code) bool) {
		for i := range DoubleBytes(x) {
			if !yield(enc.dictionary.Code(i)) {
				return
			}
		}

		if len(x)%2 == 1 {
			yield(word.Backspace)
		}
	}
}

func (enc Encoder[T]) Encode(x []byte) word.Word {
	ret := make(word.Word, 0, 1+len(x)/2)
	for c := range enc.Iter(x) {
		ret = append(ret, c)
	}
	return ret
}

func (enc Encoder[T]) Decoder() *Decoder {
	return &Decoder{
		dictionary: enc.dictionary.Inverse(),
	}
}

type Decoder struct {
	dictionary inverseDict
}

func (dec Decoder) Decode(w word.Word) ([]byte, error) {
	ret := make([]byte, 0, len(w)*2)
	for _, c := range w {
		if c == word.Backspace {
			if len(ret) == 0 {
				return nil, nil
			}
			return ret[:len(ret)-1], nil
		}
		val := dec.dictionary.Code(c)
		ret = append(ret, byte(val), byte(val>>8))
	}
	return ret, nil
}
