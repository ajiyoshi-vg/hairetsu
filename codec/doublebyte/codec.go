package doublebyte

import (
	"iter"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

var (
	_ codec.Encoder[[]byte] = (*Encoder[MapDict])(nil)
	_ codec.Decoder[[]byte] = (*Decoder)(nil)
)

type Encoder[D Dict] struct {
	dict D
}

func NewEncoder[D Dict](dict D) Encoder[D] {
	ret := Encoder[D]{
		dict: dict,
	}
	return ret
}

func DoubleBytes(x []byte) iter.Seq[uint16] {
	return func(yield func(uint16) bool) {
		for i := 0; i < len(x); i += 2 {
			var val uint16
			for j := 0; j < 2; j++ {
				if i+j < len(x) {
					val |= uint16(x[i+j]) << (8 * uint(j))
				}
			}
			if !yield(val) {
				return
			}
		}
	}
}

func (enc Encoder[D]) Iter(x []byte) iter.Seq[word.Code] {
	return func(yield func(word.Code) bool) {
		for i := range DoubleBytes(x) {
			if !yield(enc.dict.Code(i)) {
				return
			}
		}

		if len(x)%2 == 1 {
			yield(word.Backspace)
		}
	}
}

func (enc Encoder[D]) Encode(x []byte) word.Word {
	ret := make(word.Word, 0, 1+len(x)/2)
	for c := range enc.Iter(x) {
		ret = append(ret, c)
	}
	return ret
}

func (enc Encoder[D]) Decoder() *Decoder {
	return &Decoder{
		dict: enc.dict.Inverse(),
	}
}

type Decoder struct {
	dict inverseDict
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
		val := dec.dict.Code(c)
		ret = append(ret, byte(val), byte(val>>8))
	}
	return ret, nil
}
