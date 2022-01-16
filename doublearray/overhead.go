package doublearray

import (
	"fmt"
	"io"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Nodes interface {
	at(node.Index) (node.Node, error)
	io.WriterTo
}

var (
	_ Nodes = (*DoubleArray)(nil)
	_ Nodes = (*Mmap)(nil)
)

func ExactMatchSearchPointer(da *DoubleArray, cs word.Word) (node.Index, error) {
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
	data, err := da.at(n.GetOffset().Forward(word.EOS))
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}

func ExactMatchSearchInterface(da Nodes, cs word.Word) (node.Index, error) {
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
	data, err := da.at(n.GetOffset().Forward(word.EOS))
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}
