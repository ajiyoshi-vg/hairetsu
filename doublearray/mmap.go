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
	var index node.Index
	for _, c := range cs {
		next, err := da.Traverse(index, c)
		if err != nil {
			return 0, err
		}
		index = next
	}
	return da.getValue(index)
}

func (da *Mmap) CommonPrefixSearch(cs word.Word) ([]node.Index, error) {
	var ret []node.Index
	var index node.Index

	for _, c := range cs {
		next, err := da.Traverse(index, c)
		if err != nil {
			return ret, nil
		}
		index = next

		if data, err := da.getValue(index); err == nil {
			ret = append(ret, data)
		}
	}
	return ret, nil
}

func (da *Mmap) Traverse(parent node.Index, c word.Code) (node.Index, error) {
	p, err := da.at(parent)
	if err != nil {
		return 0, err
	}

	child := p.GetOffset().Forward(c)
	if int(child) >= da.length() {
		return 0, fmt.Errorf("Traverse(%d, %d) failed : index out of range", parent, c)
	}

	ch, err := da.at(child)
	if err != nil {
		return 0, err
	}

	if parent != ch.GetParent() {
		return 0, errNotChild
	}
	return child, nil
}

func (da *Mmap) getValue(index node.Index) (node.Index, error) {
	n, err := da.at(index)
	if err != nil {
		return 0, err
	}
	if !n.IsTerminal() {
		return 0, errors.New("not a terminal")
	}
	offset, err := da.Traverse(index, word.EOS)
	if err != nil {
		return 0, err
	}
	data, err := da.at(offset)
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}
