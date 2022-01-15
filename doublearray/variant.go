package doublearray

import (
	"fmt"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/runedict"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Bytes struct{}

func (Bytes) ExactMatchSearch(da *DoubleArray, cs []byte) (node.Index, error) {
	var index node.Index
	n, err := da.at(index)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		next := n.GetOffset().Forward(word.Code(c))
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
	return da.getValue(n)
}

func (Bytes) CommonPrefixSearch(da *DoubleArray, cs []byte) ([]node.Index, error) {
	var ret []node.Index
	var index node.Index
	n, err := da.at(index)
	if err != nil {
		return nil, err
	}

	for _, c := range cs {
		next := n.GetOffset().Forward(word.Code(c))
		n, err = da.at(next)
		if err != nil {
			return ret, err
		}
		if !n.IsChildOf(index) {
			return ret, nil
		}
		index = next
		if n.IsTerminal() {
			if data, err := da.getValue(n); err == nil {
				ret = append(ret, data)
			}
		}
	}
	return ret, nil
}

type Strings runedict.RuneDict

func (s Strings) ExactMatchSearch(da *DoubleArray, cs string) (node.Index, error) {
	var index node.Index
	n, err := da.at(index)
	if err != nil {
		return 0, err
	}
	for _, c := range cs {
		next := n.GetOffset().Forward(s[c])
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
	return da.getValue(n)
}

func (s Strings) CommonPrefixSearch(da *DoubleArray, cs string) ([]node.Index, error) {
	var ret []node.Index
	var index node.Index
	n, err := da.at(index)
	if err != nil {
		return nil, err
	}

	for _, c := range cs {
		next := n.GetOffset().Forward(s[c])
		n, err = da.at(next)
		if err != nil {
			return ret, err
		}
		if !n.IsChildOf(index) {
			return ret, nil
		}
		index = next
		if n.IsTerminal() {
			if data, err := da.getValue(n); err == nil {
				ret = append(ret, data)
			}
		}
	}
	return ret, nil
}
