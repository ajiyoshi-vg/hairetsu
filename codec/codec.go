package codec

import (
	"io"
	"iter"

	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Encoder[Item any] interface {
	Iter(Item) iter.Seq[word.Code]
	Encode(Item) word.Word
}
type Decoder[Item any] interface {
	Decode(word.Word) (Item, error)
}

type Dict[T, Val any] interface {
	Code(T) Val
	Inverse() Dict[Val, T]
}
type Fillable[Unit comparable] interface {
	Fill(map[Unit]int)
}
type WordDict[T comparable] interface {
	Dict[T, word.Code]
	Fillable[T]
	io.WriterTo
	io.ReaderFrom
}
