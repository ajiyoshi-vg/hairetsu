package doublearray

import (
	"fmt"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Nodes interface {
	at(node.Index) (node.Node, error)
	length() int
	getValue(node.Node) (node.Index, error)
}

func ExactMatchSearchPointer(da *DoubleArray, cs word.Word) (node.Index, error) {
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
	return da.getValue(n)
}

func ExactMatchSearchInterface(da Nodes, cs word.Word) (node.Index, error) {
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
	return da.getValue(n)
}

/*
// need go 1.18
func ExactMatchSearchGenerics[T Nodes](da T, cs word.Word) (node.Index, error) {
	var index node.Index
	length := node.Index(da.length())

	for _, c := range cs {
		next := da.at(index).GetOffset().Forward(c)
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
*/
