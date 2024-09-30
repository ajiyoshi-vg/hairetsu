package bytes

import (
	"io"
	"iter"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/codec/dict"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Dict interface {
	codec.Dict[byte, word.Code]
	io.WriterTo
	io.ReaderFrom
}
type WordDict codec.WordDict[byte]

func NewMapDict() dict.Map[byte]            { return dict.Map[byte]{} }
func NewArrayDict() dict.Array[byte]        { return dict.NewArray[byte]() }
func NewIdentityDict() *dict.Identity[byte] { return dict.NewIdentity[byte]() }

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

func (enc *Encoder[T]) Iter(x []byte) iter.Seq[word.Code] {
	return func(yield func(word.Code) bool) {
		for _, c := range x {
			if !yield(enc.dictionary.Code(c)) {
				return
			}
		}
	}
}

func (enc *Encoder[T]) Encode(x []byte) word.Word {
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

func (enc Encoder[T]) WriteTo(w io.Writer) (int64, error) {
	return enc.dictionary.WriteTo(w)
}
func (enc Encoder[T]) ReadFrom(r io.Reader) (int64, error) {
	return enc.dictionary.ReadFrom(r)
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
