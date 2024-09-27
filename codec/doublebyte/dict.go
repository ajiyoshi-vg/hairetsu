package doublebyte

import (
	"io"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Dict interface {
	codec.Dict[uint16, word.Code]
}

type Fillable interface {
	codec.Fillable[uint16]
}

type FillableDict interface {
	Dict
	Fillable
}

type WordDict interface {
	Dict
	Fillable
	io.WriterTo
	io.ReaderFrom
}
type inverseDict codec.Dict[word.Code, uint16]
