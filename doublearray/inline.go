package doublearray

import (
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

func InitialTarget(da Nodes) (node.Node, node.Index, error) {
	var parent node.Index
	target, err := da.At(parent)
	if err != nil {
		return 0, 0, err
	}
	return target, parent, nil
}

func NextTarget(da Nodes, c word.Code, target node.Node, parent node.Index) (node.Node, node.Index, error) {
	var err error
	child := target.GetChild(c)
	target, err = da.At(child)
	if err != nil {
		return 0, 0, err
	}
	if !target.IsChildOf(parent) {
		return 0, 0, ErrNotAChild
	}
	return target, child, nil
}
