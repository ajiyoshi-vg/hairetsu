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
	length := node.Index(da.length())

	for _, c := range cs {
		next := da.at(index).GetOffset().Forward(word.Code(c))
		if next >= length || !da.at(next).IsChildOf(index) {
			return 0, fmt.Errorf("ExactMatchSearch(%v) : error broken index", cs)
		}
		index = next
	}
	if !da.at(index).IsTerminal() {
		return 0, fmt.Errorf("ExactMatchSearch(%v) : not stored", cs)
	}
	data := da.at(index).GetOffset().Forward(word.EOS)
	if data >= length || !da.at(data).IsChildOf(index) {
		return 0, fmt.Errorf("ExactMatchSearch(%v) : error broken data node", cs)
	}
	return da.at(data).GetOffset(), nil
}

func (Bytes) CommonPrefixSearch(da *DoubleArray, cs []byte) ([]node.Index, error) {
	ret := make([]node.Index, 0, 10)
	var index node.Index
	length := node.Index(da.length())

	for _, c := range cs {
		next := da.at(index).GetOffset().Forward(word.Code(c))
		if next >= length || !da.at(next).IsChildOf(index) {
			return ret, nil
		}
		index = next

		if da.at(index).IsTerminal() {
			data := da.at(index).GetOffset().Forward(word.EOS)
			if data >= length || !da.at(data).IsChildOf(index) {
				err := fmt.Errorf("CommonPrefixSearch(%v) : error broken data node", cs)
				return nil, err
			}
			ret = append(ret, da.at(data).GetOffset())
		}
	}
	return ret, nil
}

type Strings runedict.RuneDict

func (s Strings) ExactMatchSearch(da *DoubleArray, cs string) (node.Index, error) {
	var index node.Index
	length := node.Index(da.length())

	for _, c := range cs {
		next := da.at(index).GetOffset().Forward(s[c])
		if next >= length || !da.at(next).IsChildOf(index) {
			return 0, fmt.Errorf("ExactMatchSearch(%v) : error broken index", cs)
		}
		index = next
	}
	if !da.at(index).IsTerminal() {
		return 0, fmt.Errorf("ExactMatchSearch(%v) : not stored", cs)
	}
	data := da.at(index).GetOffset().Forward(word.EOS)
	if data >= length || !da.at(data).IsChildOf(index) {
		return 0, fmt.Errorf("ExactMatchSearch(%v) : error broken data node", cs)
	}
	return da.at(data).GetOffset(), nil
}

func (s Strings) CommonPrefixSearch(da *DoubleArray, cs string) ([]node.Index, error) {
	ret := make([]node.Index, 0, 10)
	var index node.Index
	length := node.Index(da.length())

	for _, c := range cs {
		next := da.at(index).GetOffset().Forward(s[c])
		if next >= length || !da.at(next).IsChildOf(index) {
			return ret, nil
		}
		index = next

		if da.at(index).IsTerminal() {
			data := da.at(index).GetOffset().Forward(word.EOS)
			if data >= length || !da.at(data).IsChildOf(index) {
				err := fmt.Errorf("CommonPrefixSearch(%v) : error broken data node", cs)
				return nil, err
			}
			ret = append(ret, da.at(data).GetOffset())
		}
	}
	return ret, nil
}
