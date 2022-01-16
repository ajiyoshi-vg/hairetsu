package doublearray

import (
	"fmt"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Words struct{}

func (Words) ExactMatchSearch(da *DoubleArray, cs word.Word) (node.Index, error) {
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

func (Words) CommonPrefixSearch(da *DoubleArray, cs word.Word) ([]node.Index, error) {
	var ret []node.Index
	var index node.Index
	n, err := da.at(index)
	if err != nil {
		return ret, nil
	}

	for _, c := range cs {
		next := n.GetOffset().Forward(c)
		n, err = da.at(next)
		if err != nil {
			return ret, nil
		}
		if !n.IsChildOf(index) {
			return ret, nil
		}
		index = next
		if n.IsTerminal() {
			data, err := da.at(n.GetOffset().Forward(word.EOS))
			if err != nil {
				return ret, nil
			}
			ret = append(ret, data.GetOffset())
		}
	}
	return ret, nil
}
