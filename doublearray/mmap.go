package doublearray

import (
	"io"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/pkg/errors"
	"golang.org/x/exp/mmap"
)

type Mmap struct {
	r      *mmap.ReaderAt
	offset int64
}

const nodeSize = 8

func NewMmap(r *mmap.ReaderAt, offset int64) *Mmap {
	return &Mmap{r: r, offset: offset}
}

func OpenMmap(path string) (*Mmap, error) {
	// will be closed by runtime.SetFinalizer
	r, err := mmap.Open(path)
	if err != nil {
		return nil, err
	}
	return NewMmap(r, 0), nil
}

func (da *Mmap) At(i node.Index) (node.Node, error) {
	s := make([]byte, nodeSize)
	n, err := da.r.ReadAt(s, int64(i)*nodeSize+da.offset)
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

func (da *Mmap) ExactMatchSearch(cs word.Word) (node.Index, error) {
	return Words{}.ExactMatchSearch(da, cs)
}

func (da *Mmap) CommonPrefixSearch(cs word.Word) ([]node.Index, error) {
	return Words{}.CommonPrefixSearch(da, cs)
}

func (da *Mmap) WriteTo(w io.Writer) (int64, error) {
	var ret int64
	for i := 0; ; i++ {
		nod, err := da.At(node.Index(i))
		if err != nil {
			return ret, nil
		}
		buf, err := nod.MarshalBinary()
		if err != nil {
			return ret, err
		}
		n, err := w.Write(buf)
		ret += int64(n)
		if err != nil {
			return ret, err
		}
	}
}
