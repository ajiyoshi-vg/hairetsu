package bytes

import (
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type Walker interface {
	Walk(func([]byte) error) error
}

type Factory interface {
	Put(item.Item)
}

func FromSlice(xs [][]byte, f Factory) (Dict, error) {
	b := NewBuilder()
	for _, x := range xs {
		b.Add(x)
	}
	dict := b.Build()

	for i, x := range xs {
		w := dict.Word(x)
		f.Put(item.New(w, uint32(i)))
	}
	return dict, nil
}
