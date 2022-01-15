package doublearray

import (
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/pkg/errors"
	"golang.org/x/exp/mmap"
)

type Mmap struct {
	r *mmap.ReaderAt
}

const nodeSize = 8

func NewMmap(path string) (*Mmap, error) {
	// will be closed by runtime.SetFinalizer
	r, err := mmap.Open(path)
	if err != nil {
		return nil, err
	}
	return &Mmap{r: r}, nil
}

func (da *Mmap) at(i node.Index) (node.Node, error) {
	s := make([]byte, nodeSize)
	n, err := da.r.ReadAt(s, int64(i)*nodeSize)
	if n != nodeSize {
		return 0, errors.New("bad size")
	}
	if err != nil {
		return 0, err
	}
	var ret node.Node
	if err := ret.UnmarshalBinary(s); err != nil {
		return 0, err
	}
	return ret, nil
}

func (da *Mmap) length() int {
	return da.r.Len()
}

func (da *Mmap) ExactMatchSearch(cs word.Word) (node.Index, error) {
	return WordsMmap{}.ExactMatchSearch(da, cs)
}

func (da *Mmap) CommonPrefixSearch(cs word.Word) ([]node.Index, error) {
	return WordsMmap{}.CommonPrefixSearch(da, cs)
}
