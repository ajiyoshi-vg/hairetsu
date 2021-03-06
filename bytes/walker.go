package bytes

import (
	"github.com/ajiyoshi-vg/hairetsu/keytree"
)

type Walker interface {
	Walk(func([]byte) error) error
}

func FromWalker(w Walker) (*keytree.Tree, Dict, error) {
	b := NewBuilder()
	w.Walk(func(bs []byte) error {
		b.Add(bs)
		return nil
	})
	dict := b.Build()

	ks := keytree.New()
	var i uint32
	err := w.Walk(func(bs []byte) error {
		defer func() { i++ }()
		w := dict.Word(bs)
		return ks.Put(w, i)
	})
	if err != nil {
		return nil, nil, err
	}
	return ks, dict, err
}

func FromSlice(xs [][]byte) (*keytree.Tree, Dict, error) {
	b := NewBuilder()
	for _, x := range xs {
		b.Add(x)
	}
	dict := b.Build()

	ks := keytree.New()
	for i, x := range xs {
		if err := ks.Put(dict.Word(x), uint32(i)); err != nil {
			return nil, nil, err
		}
	}
	return ks, dict, nil
}
