package codec

import (
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

type Dict[Key, Val any] interface {
	Code(Key) Val
	Inverse() Dict[Val, Key]
}
type Fillable[Key comparable] interface {
	Fill(map[Key]int)
}
