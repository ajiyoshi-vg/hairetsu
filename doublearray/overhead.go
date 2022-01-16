package doublearray

import (
	"io"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Nodes interface {
	At(node.Index) (node.Node, error)
	io.WriterTo
}

var (
	_ Nodes = (*DoubleArray)(nil)
	_ Nodes = (*Mmap)(nil)
)

func ExactMatchSearchPointer(da *DoubleArray, cs word.Word) (node.Index, error) {
	var index node.Index
	nod, err := da.At(index)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		next := nod.GetOffset().Forward(c)
		nod, err = da.At(next)
		if err != nil {
			return 0, err
		}
		if !nod.IsChildOf(index) {
			return 0, ErrNotAChild
		}
		index = next
	}
	if !nod.IsTerminal() {
		return 0, ErrNotATerminal
	}
	data, err := da.At(nod.GetOffset().Forward(word.EOS))
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}

func ExactMatchSearchInterface(da Nodes, cs word.Word) (node.Index, error) {
	var index node.Index
	nod, err := da.At(index)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		next := nod.GetOffset().Forward(c)
		nod, err = da.At(next)
		if err != nil {
			return 0, err
		}
		if !nod.IsChildOf(index) {
			return 0, ErrNotAChild
		}
		index = next
	}
	if !nod.IsTerminal() {
		return 0, ErrNotATerminal
	}
	data, err := da.At(nod.GetOffset().Forward(word.EOS))
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}

func ExactMatchSearchPointerMmap(da *Mmap, cs word.Word) (node.Index, error) {
	var index node.Index
	nod, err := da.At(index)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		next := nod.GetOffset().Forward(c)
		nod, err = da.At(next)
		if err != nil {
			return 0, err
		}
		if !nod.IsChildOf(index) {
			return 0, ErrNotAChild
		}
		index = next
	}
	if !nod.IsTerminal() {
		return 0, ErrNotATerminal
	}
	data, err := da.At(nod.GetOffset().Forward(word.EOS))
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}
