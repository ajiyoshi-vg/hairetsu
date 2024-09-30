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
	target, parent, err := doublearray.InitialTarget(da)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		target, parent, err = doublearray.NextTarget(da, word.Code(c), target, parent)
		if err != nil {
			return 0, err
		}
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

func ExactMatchSearchInterface(da doublearray.Nodes, cs []byte) (node.Index, error) {
	target, parent, err := doublearray.InitialTarget(da)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		target, parent, err = doublearray.NextTarget(da, word.Code(c), target, parent)
		if err != nil {
			return 0, err
		}
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
	target, parent, err := doublearray.InitialTarget(da)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		target, parent, err = doublearray.NextTarget(da, word.Code(c), target, parent)
		if err != nil {
			return 0, err
		}
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

func ExactMatchSearchGenerics[T doublearray.Nodes](da T, cs []byte) (node.Index, error) {
	target, parent, err := doublearray.InitialTarget(da)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		target, parent, err = doublearray.NextTarget(da, word.Code(c), target, parent)
		if err != nil {
			return 0, err
		}
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
