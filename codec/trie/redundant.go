package trie

import (
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

func BytesExactMatchSearch[DA doublearray.Nodes](da DA, cs []byte) (node.Index, error) {
	target, parent, err := doublearray.InitialTarget(da)
	if err != nil {
		return 0, err
	}

	for _, c := range cs {
		code := word.Code(c)
		target, parent, err = doublearray.NextTarget(da, code, target, parent)
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

func BytesCommonPrefixSearch[DA doublearray.Nodes](da DA, cs []byte) ([]node.Index, error) {
	var ret []node.Index
	target, parent, err := doublearray.InitialTarget(da)
	if err != nil {
		return nil, err
	}

	for _, c := range cs {
		code := word.Code(c)
		target, parent, err = doublearray.NextTarget(da, code, target, parent)
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

func RunsExactMatchSearch[DA doublearray.Nodes](da DA, cs string) (node.Index, error) {
	target, parent, err := doublearray.InitialTarget(da)
	if err != nil {
		return 0, err
	}

	for _, c := range cs {
		code := word.Code(c)
		target, parent, err = doublearray.NextTarget(da, code, target, parent)
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

func RunesCommonPrefixSearch[DA doublearray.Nodes](da DA, cs string) ([]node.Index, error) {
	var ret []node.Index
	target, parent, err := doublearray.InitialTarget(da)
	if err != nil {
		return nil, err
	}

	for _, c := range cs {
		code := word.Code(c)
		target, parent, err = doublearray.NextTarget(da, code, target, parent)
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
