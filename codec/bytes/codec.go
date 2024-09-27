package bytes

import (
	"iter"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Dict codec.Dict[byte, word.Code]
type WordDict codec.WordDict[byte]

var (
	_ codec.Encoder[[]byte] = (*Encoder[Dict])(nil)
	_ codec.Decoder[[]byte] = (*Decoder)(nil)
)

type Encoder[T Dict] struct {
	dictionary T
}

func NewEncoder[D Dict](d D) *Encoder[D] {
	return &Encoder[D]{
		dictionary: d,
	}
}

func (enc Encoder[T]) Iter(x []byte) iter.Seq[word.Code] {
	return func(yield func(word.Code) bool) {
		for _, c := range x {
			if !yield(enc.dictionary.Code(c)) {
				return
			}
		}
	}
}

func (enc Encoder[T]) Encode(x []byte) word.Word {
	ret := make(word.Word, 0, len(x))
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
	dictionary codec.Dict[word.Code, byte]
}

func (dec *Decoder) Decode(x word.Word) ([]byte, error) {
	ret := make([]byte, 0, len(x))
	for _, c := range x {
		b := dec.dictionary.Code(c)
		ret = append(ret, b)
	}
	return ret, nil
}
