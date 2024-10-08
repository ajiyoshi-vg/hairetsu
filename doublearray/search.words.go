package doublearray

import (
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Words struct{}

func (Words) ExactMatchSearch(da Nodes, cs word.Word) (node.Index, error) {
	target, parent, err := InitialTarget(da)
	if err != nil {
		return 0, err
	}

	for _, c := range cs {
		code := c
		target, parent, err = NextTarget(da, code, target, parent)
		if err != nil {
			return 0, err
		}
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
	target, parent, err := InitialTarget(da)
	if err != nil {
		return nil, err
	}

	for _, c := range cs {
		code := c
		target, parent, err = NextTarget(da, code, target, parent)
		if err != nil {
			return ret, nil
		}

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
