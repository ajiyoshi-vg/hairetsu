package trie

import (
	"io"
	"os"

	"github.com/ajiyoshi-vg/hairetsu/codec/u16s"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
)

type Inline[D u16s.WordDict] struct {
	dict D
	data *doublearray.DoubleArray
}

func NewInline[D u16s.WordDict](dict D) *Inline[D] {
	return &Inline[D]{dict: dict}
}

func OpenInline[D u16s.WordDict](path string, dict D) (*Inline[D], error) {
	ret := NewInline(dict)
	if err := ret.Open(path); err != nil {
		return nil, err
	}
	return ret, nil
}

func (x *Inline[D]) ReadFrom(r io.Reader) (int64, error) {
	var ret int64
	n, err := x.dict.ReadFrom(r)
	ret += n
	if err != nil {
		return ret, err
	}
	da := doublearray.New()
	m, err := da.ReadFrom(r)
	ret += m
	if err != nil {
		return ret, err
	}
	x.data = da
	return ret, nil
}

func (x *Inline[D]) Open(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := x.ReadFrom(file); err != nil {
		return err
	}
	return nil
}

func (x *Inline[D]) Searcher() *u16s.InlineSearcher[D, *doublearray.DoubleArray] {
	return u16s.NewInlineSearcher(x.data, x.dict)
}
