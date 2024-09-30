package trie

import (
	"bufio"
	"io"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"golang.org/x/exp/mmap"
)

type Mmap[X any] struct {
	enc codec.Encoder[X]
	da  *doublearray.Mmap
	m   *mmap.ReaderAt
}

func NewMmap[X any](enc codec.Encoder[X]) *Mmap[X] {
	return &Mmap[X]{enc: enc}
}

func OpenMmap[X any](path string, enc codec.Encoder[X]) (*Mmap[X], error) {
	ret := NewMmap(enc)
	if err := ret.Open(path); err != nil {
		return nil, err
	}
	return ret, nil
}

func (t *Mmap[X]) Open(path string) error {
	m, err := mmap.Open(path)
	if err != nil {
		return err
	}
	if err := t.open(m); err != nil {
		m.Close()
		return err
	}
	t.m = m
	return nil
}

func (t *Mmap[X]) open(r *mmap.ReaderAt) error {
	sec := io.NewSectionReader(r, 0, int64(r.Len()))
	br := bufio.NewReader(sec)
	n, err := t.enc.ReadFrom(br)
	if err != nil {
		return err
	}

	t.da = doublearray.NewMmap(r, int64(n), int64(r.Len())-int64(n))
	return nil
}

func (t *Mmap[X]) Close() error {
	return t.m.Close()
}

func (t *Mmap[X]) Searcher() *Searcher[X, *doublearray.Mmap] {
	return NewSearcher(t.enc, t.da)
}

func (t *Mmap[X]) WriteTo(w io.Writer) (int64, error) {
	return multiCopy(w, t.enc, t.da)
}
