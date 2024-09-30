package composer

import (
	"bufio"
	"io"
	"os"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"golang.org/x/exp/mmap"
)

type FileTrie[X any] struct {
	enc codec.Encoder[X]
	da  *doublearray.DoubleArray
}

func NewFileTrie[X any](enc codec.Encoder[X], opt ...Option[FileTrie[X]]) *FileTrie[X] {
	ret := &FileTrie[X]{enc: enc}
	for _, f := range opt {
		f(ret)
	}
	return ret
}

func (t *FileTrie[X]) Searcher() *codec.Searcher[X, *doublearray.DoubleArray] {
	return codec.NewSearcher(t.enc, t.da)
}

func (t *FileTrie[X]) WriteTo(w io.Writer) (int64, error) {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	var ret int64
	{ // write encoder
		n, err := t.enc.WriteTo(bw)
		ret += n
		if err != nil {
			return ret, err
		}
	}
	{ // write index
		n, err := t.da.WriteTo(bw)
		ret += n
		if err != nil {
			return ret, err
		}
	}
	return 0, nil
}

func (t *FileTrie[X]) ReadFrom(r io.Reader) (int64, error) {
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

func (t *FileTrie[X]) Open(path string) error {
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

type MmapTrie[X any] struct {
	enc codec.Encoder[X]
	da  *doublearray.Mmap
}

func NewMmapTrie[X any](enc codec.Encoder[X]) *MmapTrie[X] {
	return &MmapTrie[X]{enc: enc}
}

var _ io.ReaderAt

func (t *MmapTrie[X]) Open(path string) error {
	m, err := mmap.Open(path)
	if err != nil {
		return err
	}
	if err := t.open(m); err != nil {
		m.Close()
		return err
	}
	return nil
}

func (t *MmapTrie[X]) open(r *mmap.ReaderAt) error {
	sec := io.NewSectionReader(r, 0, int64(r.Len()))
	br := bufio.NewReader(sec)
	n, err := t.enc.ReadFrom(br)
	if err != nil {
		return err
	}

	t.da = doublearray.NewMmap(r, int64(n), int64(r.Len())-int64(n))
	return nil
}

func (t *MmapTrie[X]) Searcher() *codec.Searcher[X, *doublearray.Mmap] {
	return codec.NewSearcher(t.enc, t.da)
}

func (t *MmapTrie[X]) WriteTo(w io.Writer) (int64, error) {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	var ret int64
	{ // write encoder
		n, err := t.enc.WriteTo(bw)
		ret += n
		if err != nil {
			return ret, err
		}
	}
	{ // write index
		n, err := t.da.WriteTo(bw)
		ret += n
		if err != nil {
			return ret, err
		}
	}
	return 0, nil
}
