package bytes

import "github.com/ajiyoshi-vg/hairetsu/keytree"

type Walker interface {
	Walk(func([]byte) error) error
}

func FromBytesWalker(w Walker) (*keytree.Tree, Dict, error) {
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
