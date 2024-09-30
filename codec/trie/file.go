package trie

import (
	"bufio"
	"io"
	"os"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
)

type File[X any] struct {
	enc codec.Encoder[X]
	da  *doublearray.DoubleArray
}

func NewFile[X any](enc codec.Encoder[X], opt ...Option[File[X]]) *File[X] {
	ret := &File[X]{enc: enc}
	for _, f := range opt {
		f(ret)
	}
	return ret
}

func OpenFile[X any](path string, enc codec.Encoder[X]) (*File[X], error) {
	ret := NewFile(enc)
	if err := ret.Open(path); err != nil {
		return nil, err
	}
	return ret, nil
}

func (t *File[X]) Searcher() *Searcher[X, *doublearray.DoubleArray] {
	return NewSearcher(t.enc, t.da)
}

func (t *File[X]) Index() *doublearray.DoubleArray {
	return t.da
}

func (t *File[X]) WriteTo(w io.Writer) (int64, error) {
	return multiCopy(w, t.enc, t.da)
}

func (t *File[X]) ReadFrom(r io.Reader) (int64, error) {
	br := bufio.NewReader(r)
	var ret int64
	{ // read encoder
		n, err := t.enc.ReadFrom(br)
		ret += n
		if err != nil {
			return ret, err
		}
	}
	{ // read index
		da := doublearray.New()
		n, err := da.ReadFrom(br)
		ret += n
		if err != nil {
			return ret, err
		}
		t.da = da
	}
	return 0, nil
}

func (t *File[X]) Open(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := t.ReadFrom(f); err != nil {
		return err
	}
	return nil
}

func (t *File[X]) Stat() doublearray.Stat {
	return doublearray.GetStat(t.da)
}
