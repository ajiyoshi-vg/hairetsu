package runes

import (
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type Factory interface {
	Put(item.Item)
}

func FromSlice(xs []string, f Factory) (Dict, error) {
	b := NewBuilder()
	for _, x := range xs {
		b.Add(x)
	}
	dict := b.Build()

	for i, x := range xs {
		w, err := dict.Word(x)
		if err != nil {
			return nil, err
		}
		f.Put(item.New(w, uint32(i)))
	}
	return dict, nil
}
