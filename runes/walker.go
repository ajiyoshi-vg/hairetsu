package runes

import "github.com/ajiyoshi-vg/hairetsu/keytree"

type Walker interface {
	Walk(func(string) error) error
}

func FromWalker(w Walker) (*keytree.Tree, Dict, error) {
	b := NewBuilder()
	w.Walk(func(bs string) error {
		b.Add(bs)
		return nil
	})
	dict := b.Build()

	ks := keytree.New()
	var i uint32
	err := w.Walk(func(bs string) error {
		defer func() { i++ }()
		w, err := dict.Word(bs)
		if err != nil {
			return err
		}
		return ks.Put(w, i)
	})
	if err != nil {
		return nil, nil, err
	}
	return ks, dict, err
}

func FromSlice(xs []string) (*keytree.Tree, Dict, error) {
	b := NewBuilder()
	for _, x := range xs {
		b.Add(x)
	}
	dict := b.Build()

	ks := keytree.New()
	for i, x := range xs {
		w, err := dict.Word(x)
		if err != nil {
			return nil, nil, err
		}
		if err := ks.Put(w, uint32(i)); err != nil {
			return nil, nil, err
		}
	}
	return ks, dict, nil
}
