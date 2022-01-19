package overhead

import (
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Nodes interface {
	At(node.Index) (node.Node, error)
}

var (
	_ Nodes = (*doublearray.DoubleArray)(nil)
	_ Nodes = (*doublearray.Mmap)(nil)
)

func ExactMatchSearchPointer(da *doublearray.DoubleArray, cs []byte) (node.Index, error) {
	var parent node.Index
	target, err := da.At(parent)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		child := target.GetChild(word.Code(c))
		target, err = da.At(child)
		if err != nil {
			return 0, err
		}
		if !target.IsChildOf(parent) {
			return 0, doublearray.ErrNotAChild
		}
		parent = child
	}
	if !target.IsTerminal() {
		return 0, doublearray.ErrNotATerminal
	}
	data, err := da.At(target.GetChild(word.EOS))
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}

func ExactMatchSearchInterface(da Nodes, cs []byte) (node.Index, error) {
	var parent node.Index
	target, err := da.At(parent)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		child := target.GetChild(word.Code(c))
		target, err = da.At(child)
		if err != nil {
			return 0, err
		}
		if !target.IsChildOf(parent) {
			return 0, doublearray.ErrNotAChild
		}
		parent = child
	}
	if !target.IsTerminal() {
		return 0, doublearray.ErrNotATerminal
	}
	data, err := da.At(target.GetChild(word.EOS))
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}

func ExactMatchSearchPointerMmap(da *doublearray.Mmap, cs []byte) (node.Index, error) {
	var parent node.Index
	target, err := da.At(parent)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		child := target.GetChild(word.Code(c))
		target, err = da.At(child)
		if err != nil {
			return 0, err
		}
		if !target.IsChildOf(parent) {
			return 0, doublearray.ErrNotAChild
		}
		parent = child
	}
	if !target.IsTerminal() {
		return 0, doublearray.ErrNotATerminal
	}
	data, err := da.At(target.GetChild(word.EOS))
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}
