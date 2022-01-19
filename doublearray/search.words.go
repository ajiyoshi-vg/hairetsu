package doublearray

import (
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Words struct{}

func (Words) ExactMatchSearch(da Nodes, cs word.Word) (node.Index, error) {
	var parent node.Index
	target, err := da.At(parent)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		child := target.GetChild(c)
		target, err = da.At(child)
		if err != nil {
			return 0, err
		}
		if !target.IsChildOf(parent) {
			return 0, ErrNotAChild
		}
		parent = child
	}
	if !target.IsTerminal() {
		return 0, ErrNotATerminal
	}
	data, err := da.At(target.GetChild(word.EOS))
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}

func (Words) CommonPrefixSearch(da Nodes, cs word.Word) ([]node.Index, error) {
	var ret []node.Index
	var parent node.Index
	target, err := da.At(parent)
	if err != nil {
		return ret, nil
	}

	for _, c := range cs {
		child := target.GetChild(c)
		target, err = da.At(child)
		if err != nil {
			return ret, nil
		}
		if !target.IsChildOf(parent) {
			return ret, nil
		}
		parent = child
		if target.IsTerminal() {
			data, err := da.At(target.GetChild(word.EOS))
			if err != nil {
				return ret, nil
			}
			ret = append(ret, data.GetOffset())
		}
	}
	return ret, nil
}
