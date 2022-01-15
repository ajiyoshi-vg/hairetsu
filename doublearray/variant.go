package doublearray

import (
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/runedict"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Bytes struct{}

func (Bytes) ExactMatchSearch(da *DoubleArray, cs []byte) (node.Index, error) {
	var index node.Index
	for _, c := range cs {
		next, err := da.Traverse(index, word.Code(c))
		if err != nil {
			return 0, err
		}
		index = next
	}
	return da.getValue(index)
}

func (Bytes) CommonPrefixSearch(da *DoubleArray, cs []byte) ([]node.Index, error) {
	var ret []node.Index
	var index node.Index

	for _, c := range cs {
		next, err := da.Traverse(index, word.Code(c))
		if err != nil {
			return ret, nil
		}
		index = next

		if data, err := da.getValue(index); err == nil {
			ret = append(ret, data)
		}
	}
	return ret, nil
}

type Strings runedict.RuneDict

func (s Strings) ExactMatchSearch(da *DoubleArray, cs string) (node.Index, error) {
	var index node.Index
	for _, c := range cs {
		next, err := da.Traverse(index, s[c])
		if err != nil {
			return 0, err
		}
		index = next
	}
	return da.getValue(index)
}

func (s Strings) CommonPrefixSearch(da *DoubleArray, cs string) ([]node.Index, error) {
	var ret []node.Index
	var index node.Index

	for _, c := range cs {
		next, err := da.Traverse(index, s[c])
		if err != nil {
			return ret, nil
		}
		index = next

		if data, err := da.getValue(index); err == nil {
			ret = append(ret, data)
		}
	}
	return ret, nil
}
