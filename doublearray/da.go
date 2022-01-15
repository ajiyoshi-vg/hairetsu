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
	return Words{}.ExactMatchSearch(da, cs)
}

func (da *DoubleArray) CommonPrefixSearch(cs word.Word) ([]node.Index, error) {
	return Words{}.CommonPrefixSearch(da, cs)
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
