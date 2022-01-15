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

func (da *DoubleArray) ExactMatchSearch(cs word.Word) (node.Index, error) {
	var index node.Index
	length := node.Index(da.length())

	for _, c := range cs {
		next := da.at(index).GetOffset().Forward(c)
		if next >= length || !da.at(next).IsChildOf(index) {
			return 0, fmt.Errorf("ExactMatchSearch(%v) : error broken index", cs)
		}
		index = next
	}
	if !da.at(index).IsTerminal() {
		return 0, fmt.Errorf("ExactMatchSearch(%v) : not stored", cs)
	}
	data := da.at(index).GetOffset().Forward(word.EOS)
	if data >= length || !da.at(data).IsChildOf(index) {
		return 0, fmt.Errorf("ExactMatchSearch(%v) : error broken data node", cs)
	}
	return da.at(data).GetOffset(), nil
}

func (da *DoubleArray) CommonPrefixSearch(cs word.Word) ([]node.Index, error) {
	ret := make([]node.Index, 0, 10)
	var index node.Index
	length := node.Index(da.length())

	for _, c := range cs {
		next := da.at(index).GetOffset().Forward(c)
		if next >= length || !da.at(next).IsChildOf(index) {
			return ret, nil
		}
		index = next

		if da.at(index).IsTerminal() {
			data := da.at(index).GetOffset().Forward(word.EOS)
			if data >= length || !da.at(data).IsChildOf(index) {
				err := fmt.Errorf("CommonPrefixSearch(%v) : error broken data node", cs)
				return nil, err
			}
			ret = append(ret, da.at(data).GetOffset())
		}
	}
	return ret, nil
}

var errNotChild = errors.New("not a child")

func (da *DoubleArray) Traverse(parent node.Index, c word.Code) (node.Index, error) {
	child := da.at(parent).GetOffset().Forward(c)
	if int(child) >= len(da.nodes) {
		return 0, fmt.Errorf("Traverse(%d, %d) failed : index out of range", parent, c)
	}
	if parent != da.at(child).GetParent() {
		return 0, errNotChild
	}
	return child, nil
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

func (da *DoubleArray) at(i node.Index) node.Node {
	return da.nodes[i]
}

func (da *DoubleArray) length() int {
	return len(da.nodes)
}
