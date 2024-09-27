package doublebyte

import (
	"io"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/codec/dict"
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

func NewMapDict() dict.MapDict[uint16] {
	return dict.MapDict[uint16]{}
}
func NewArrayDict() dict.ArrayDict[uint16] {
	return dict.NewArrayDict[uint16]()
}
func NewIdentityDict() *dict.Identity[uint16] {
	return &dict.Identity[uint16]{}
}
