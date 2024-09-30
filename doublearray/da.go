package doublearray

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

var (
	ErrNotAChild    = errors.New("not a child")
	ErrNotATerminal = errors.New("not a terminal")
)

type Nodes interface {
	At(node.Index) (node.Node, error)
	io.WriterTo
}

var (
	_ Nodes = (*DoubleArray)(nil)
	_ Nodes = (*Mmap)(nil)
)

type DoubleArray struct {
	nodes []node.Node
}

func New() *DoubleArray {
	return &DoubleArray{
		nodes: make([]node.Node, 8),
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
func (da *DoubleArray) ReadFrom(r io.Reader) (int64, error) {
	return NewBuilder().readFrom(da, r)
}

func (da *DoubleArray) ExactMatchSearch(cs word.Word) (node.Index, error) {
	return Words{}.ExactMatchSearch(da, cs)
}

func (da *DoubleArray) CommonPrefixSearch(cs word.Word) ([]node.Index, error) {
	return Words{}.CommonPrefixSearch(da, cs)
}

func (da *DoubleArray) At(i node.Index) (node.Node, error) {
	if int(i) >= len(da.nodes) {
		return 0, fmt.Errorf("index(%d) out of range", i)
	}
	return da.nodes[i], nil
}

func OpenFile(path string) (*DoubleArray, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	da := New()
	if _, err := da.ReadFrom(bufio.NewReader(file)); err != nil {
		return nil, err
	}
	return da, nil
}
