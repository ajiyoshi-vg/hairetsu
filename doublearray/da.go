package doublearray

import (
	"errors"
	"fmt"
	"io"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type DoubleArray struct {
	nodes []node.Node
}

func New() *DoubleArray {
	return &DoubleArray{
		nodes: make([]node.Node, 10),
	}
}

func FromArray(xs []uint64) *DoubleArray {
	nodes := make([]node.Node, len(xs))
	for i, x := range xs {
		nodes[i] = node.Node(x)
	}
	return &DoubleArray{
		nodes: nodes,
	}
}

func (da *DoubleArray) Array() []uint64 {
	ret := make([]uint64, len(da.nodes))
	for i, x := range da.nodes {
		ret[i] = uint64(x)
	}
	return ret
}

func (da *DoubleArray) Stat() Stat {
	return newStat(da)
}

func (da *DoubleArray) WriteTo(w io.Writer) (int64, error) {
	var ret int64
	for _, node := range da.nodes {
		buf, err := node.MarshalBinary()
		if err != nil {
			return ret, err
		}
		n, err := w.Write(buf)
		ret += int64(n)
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}

func (da *DoubleArray) ExactMatchSearch(cs word.Word) (node.Index, error) {
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

func (da *DoubleArray) CommonPrefixSearch(cs word.Word) ([]node.Index, error) {
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

var errNotChild = errors.New("not a child")

func (da *DoubleArray) Traverse(parent node.Index, c word.Code) (node.Index, error) {
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

func (da *DoubleArray) at(i node.Index) (node.Node, error) {
	if int(i) >= len(da.nodes) {
		return 0, fmt.Errorf("index(%d) out of range", i)
	}
	return da.nodes[i], nil
}

func (da *DoubleArray) length() int {
	return len(da.nodes)
}

func (da *DoubleArray) getValue(index node.Index) (node.Index, error) {
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
