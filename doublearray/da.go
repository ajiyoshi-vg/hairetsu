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

func (da *DoubleArray) CommonPrefixSearch(cs word.Word) ([]node.Index, error) {
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

var errNotChild = errors.New("not a child")

func (da *DoubleArray) at(i node.Index) (node.Node, error) {
	if int(i) >= len(da.nodes) {
		return 0, fmt.Errorf("index(%d) out of range", i)
	}
	return da.nodes[i], nil
}

func (da *DoubleArray) length() int {
	return len(da.nodes)
}

func (da *DoubleArray) getValue(n node.Node) (node.Index, error) {
	offset := n.GetOffset().Forward(word.EOS)
	data, err := da.at(offset)
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}
