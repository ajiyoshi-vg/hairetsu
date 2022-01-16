package doublearray

import (
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Words struct{}

func (Words) ExactMatchSearch(da Nodes, cs word.Word) (node.Index, error) {
	var index node.Index
	nod, err := da.at(index)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		next := nod.GetOffset().Forward(c)
		nod, err = da.at(next)
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
	data, err := da.at(nod.GetOffset().Forward(word.EOS))
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}

func (Words) CommonPrefixSearch(da Nodes, cs word.Word) ([]node.Index, error) {
	var ret []node.Index
	var index node.Index
	nod, err := da.at(index)
	if err != nil {
		return ret, nil
	}

	for _, c := range cs {
		next := nod.GetOffset().Forward(c)
		nod, err = da.at(next)
		if err != nil {
			return ret, nil
		}
		if !nod.IsChildOf(index) {
			return ret, nil
		}
		index = next
		if nod.IsTerminal() {
			data, err := da.at(nod.GetOffset().Forward(word.EOS))
			if err != nil {
				return ret, nil
			}
			ret = append(ret, data.GetOffset())
		}
	}
	return ret, nil
}
