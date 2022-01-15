package doublearray

import (
	"fmt"

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

func (da *Mmap) ExactMatchSearch(cs word.Word) (node.Index, error) {
	var index node.Index
	n, err := da.at(index)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		next := n.GetOffset().Forward(c)
		n, err = da.at(next)
		if err != nil {
			return 0, err
		}
		if !n.IsChildOf(index) {
			return 0, errNotChild
		}
		index = next
	}
	if !n.IsTerminal() {
		return 0, fmt.Errorf("not a terminal")
	}
	return da.getValue(n)
}

func (da *Mmap) CommonPrefixSearch(cs word.Word) ([]node.Index, error) {
	var ret []node.Index
	var index node.Index
	n, err := da.at(index)
	if err != nil {
		return nil, err
	}

	for _, c := range cs {
		next := n.GetOffset().Forward(c)
		n, err = da.at(next)
		if err != nil {
			return ret, err
		}
		if !n.IsChildOf(index) {
			return ret, nil
		}
		index = next
		if n.IsTerminal() {
			if data, err := da.getValue(n); err == nil {
				ret = append(ret, data)
			}
		}
	}
	return ret, nil
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

func (da *Mmap) getValue(n node.Node) (node.Index, error) {
	offset := n.GetOffset().Forward(word.EOS)
	data, err := da.at(offset)
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}
