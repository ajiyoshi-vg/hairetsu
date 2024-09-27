package runes

import (
	"iter"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/codec/dict"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Dict codec.Dict[rune, word.Code]
type WordDict codec.WordDict[rune]

func NewMapDict() dict.Map[rune]            { return dict.Map[rune]{} }
func NewIdentityDict() *dict.Identity[rune] { return dict.NewIdentity[rune]() }

var (
	_ codec.Encoder[string] = (*Encoder[Dict])(nil)
	_ codec.Decoder[string] = (*Decoder)(nil)
)

type Encoder[T Dict] struct {
	dictionary T
}

func NewEncoder[T Dict](d T) *Encoder[T] {
	return &Encoder[T]{
		dictionary: d,
	}
}

func (enc Encoder[T]) Iter(x string) iter.Seq[word.Code] {
	return func(yield func(word.Code) bool) {
		for _, c := range x {
			if !yield(enc.dictionary.Code(c)) {
				return
			}
		}
	}
}

func (enc Encoder[T]) Encode(x string) word.Word {
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
	dictionary codec.Dict[word.Code, rune]
}

func (dec *Decoder) Decode(x word.Word) (string, error) {
	ret := make([]rune, 0, len(x))
	for _, c := range x {
		b := dec.dictionary.Code(c)
		ret = append(ret, b)
	}
	return string(ret), nil
}
