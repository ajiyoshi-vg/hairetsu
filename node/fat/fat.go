package fat

import (
	"fmt"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/pkg/errors"
)

type Node struct {
	id         int
	base       node.Index
	check      node.Index
	hasParent  bool
	hasOffset  bool
	isTerminal bool
}

func Root() *Node {
	return newNode(0, 0, 1)
}
func New(i int) *Node {
	return newNode(i, node.Index(i-1), node.Index(i+1))
}
func newNode(id int, b node.Index, c node.Index) *Node {
	return &Node{id: id, base: b, check: c}
}

func (x Node) String() string {
	ret := fmt.Sprintf("{%s:%d, %s:%d}",
		x.baseLabel(),
		x.base,
		x.checkLabel(),
		x.check,
	)
	if x.IsTerminal() {
		ret += "#"
	}
	return ret
}

func (x Node) GetOffset() node.Index {
	return x.base
}
func (x *Node) SetOffset(i node.Index) {
	x.hasOffset = true
	x.setBase(i)
}
func (x *Node) Terminate() {
	x.isTerminal = true
}
func (x Node) IsTerminal() bool {
	return x.isTerminal
}
func (x Node) HasParent() bool {
	return x.hasParent
}
func (x *Node) SetParent(i node.Index) {
	x.hasParent = true
	x.setCheck(i)
}
func (x Node) IsChildOf(parent node.Index) bool {
	if !x.HasParent() {
		return false
	}
	return x.check == parent
}
func (x Node) GetNextEmptyNode() node.Index {
	if x.HasParent() {
		panic(errors.Errorf("emptyでないNode(%v)のnextEmptyNodeを取ろうとした", x))
	}
	return x.check
}
func (x *Node) SetNextEmptyNode(i node.Index) {
	x.setCheck(i)
}
func (x *Node) SetPrevEmptyNode(i node.Index) {
	x.setBase(i)
}

func (x Node) GetPrevEmptyNode() node.Index {
	if x.hasOffset {
		panic(errors.Errorf("prevが存在しないNode(%v)のGetPrevEmptyNodeを取ろうとした", x))
	}
	return x.base
}

/////

func (x *Node) setBase(y node.Index) {
	x.base = y
}
func (x *Node) setCheck(y node.Index) {
	x.check = y
}
func (x Node) baseLabel() string {
	if x.hasOffset {
		return "base"
	}
	return "prev"
}
func (x Node) checkLabel() string {
	if x.hasParent {
		return "check"
	}
	return "next"
}
