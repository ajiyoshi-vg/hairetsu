package doublearray

import (
	"github.com/ajiyoshi-vg/hairetsu/node"
)

type nodeFactory interface {
	root() node.Node
	node(int) node.Node
}

var (
	_ nodeFactory = (*factory)(nil)
)

type factory struct{}

func (f *factory) root() node.Node {
	return node.Root()
}

func (f *factory) node(i int) node.Node {
	return node.New(i)
}
