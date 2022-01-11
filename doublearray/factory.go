package doublearray

import (
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/node/fat"
)

type nodeFactory interface {
	root() node.Node
	node(int) node.Node
}

var (
	_ nodeFactory = (*fatFactory)(nil)
	_ nodeFactory = (*u64Factory)(nil)
)

type fatFactory struct{}

func (f *fatFactory) root() node.Node {
	return fat.Root()
}

func (f *fatFactory) node(i int) node.Node {
	return fat.New(i)
}

type u64Factory struct{}

func (f *u64Factory) root() node.Node {
	return fat.Root()
}

func (f *u64Factory) node(i int) node.Node {
	return fat.New(i)
}
