package hairetsu

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
)

type fatFactory struct{}

func (f *fatFactory) root() node.Node {
	return fat.Root()
}

func (f *fatFactory) node(i int) node.Node {
	return fat.New(i)
}
