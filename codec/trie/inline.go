package trie

import (
	"bufio"
	"io"
	"os"

	"github.com/ajiyoshi-vg/hairetsu/codec/u16s"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"golang.org/x/exp/mmap"
)

type InlineFile[D u16s.WordDict] struct {
	dict D
	data *doublearray.DoubleArray
}

func NewInlineFile[D u16s.WordDict](dict D) *InlineFile[D] {
	return &InlineFile[D]{dict: dict}
}

func OpenInlineFile[D u16s.WordDict](path string, dict D) (*InlineFile[D], error) {
	ret := NewInlineFile(dict)
	if err := ret.Open(path); err != nil {
		return nil, err
	}
	return ret, nil
}

func (x *InlineFile[D]) ReadFrom(r io.Reader) (int64, error) {
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

func (x *InlineFile[D]) Open(path string) error {
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

func (x *InlineFile[D]) Searcher() *u16s.InlineSearcher[D, *doublearray.DoubleArray] {
	return u16s.NewInlineSearcher(x.data, x.dict)
}

type InlineMmap[D u16s.WordDict] struct {
	dict D
	data *doublearray.Mmap
	m    *mmap.ReaderAt
}

func OpenInlineMmap[D u16s.WordDict](path string, dict D) (*InlineMmap[D], error) {
	ret := NewInlineMmap(dict)
	if err := ret.Open(path); err != nil {
		return nil, err
	}
	return ret, nil
}

func NewInlineMmap[D u16s.WordDict](dict D) *InlineMmap[D] {
	return &InlineMmap[D]{dict: dict}
}

func (x *InlineMmap[D]) Open(path string) error {
	m, err := mmap.Open(path)
	if err != nil {
		return err
	}
	if err := x.open(m); err != nil {
		m.Close()
		return err
	}
	x.m = m
	return nil
}

func (x *InlineMmap[D]) open(r *mmap.ReaderAt) error {
	sec := io.NewSectionReader(r, 0, int64(r.Len()))
	br := bufio.NewReader(sec)
	n, err := x.dict.ReadFrom(br)
	if err != nil {
		return err
	}
	x.data = doublearray.NewMmap(r, int64(n), int64(r.Len())-int64(n))
	return nil
}

func (x *InlineMmap[D]) Close() error {
	return x.m.Close()
}

func (x *InlineMmap[D]) Searcher() *u16s.InlineSearcher[D, *doublearray.Mmap] {
	return u16s.NewInlineSearcher(x.data, x.dict)
}
