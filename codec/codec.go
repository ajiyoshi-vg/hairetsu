package codec

import (
	"io"
	"iter"

	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Encoder[X any] interface {
	Iter(X) iter.Seq[word.Code]
	Encode(X) word.Word
}
type Decoder[X any] interface {
	Decode(word.Word) (X, error)
}

type Dict[T, Val any] interface {
	Code(T) Val
	Inverse() Dict[Val, T]
}
type Fillable[T comparable] interface {
	Fill(map[T]int)
	io.WriterTo
	io.ReaderFrom
}
type FillableDict[T comparable] interface {
	Dict[T, word.Code]
	Fillable[T]
}
type WordDict[T comparable] interface {
	Dict[T, word.Code]
	Fillable[T]
}
