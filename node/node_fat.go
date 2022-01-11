// +build verbose

package node

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

type Node struct {
	id         int
	base       Index
	check      Index
	hasParent  bool
	hasOffset  bool
	isTerminal bool
}

var _ NodeInterface = (*Node)(nil)

func Root() Node {
	return newNode(0, 0, 1)
}
func New(i int) Node {
	return newNode(i, Index(i-1), Index(i+1))
}
func newNode(id int, b Index, c Index) Node {
	return Node{id: id, base: b, check: c}
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

func (x Node) GetOffset() Index {
	return x.base
}
func (x *Node) SetOffset(i Index) {
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
func (x *Node) SetParent(i Index) {
	x.hasParent = true
	x.setCheck(i)
}
func (x Node) IsChildOf(parent Index) bool {
	if !x.HasParent() {
		return false
	}
	return x.check == parent
}
func (x Node) GetNextEmptyNode() Index {
	if x.HasParent() {
		panic(errors.Errorf("emptyでないNode(%v)のnextEmptyNodeを取ろうとした", x))
	}
	return x.check
}
func (x *Node) SetNextEmptyNode(i Index) {
	x.setCheck(i)
}
func (x *Node) SetPrevEmptyNode(i Index) {
	x.setBase(i)
}

func (x Node) GetPrevEmptyNode() Index {
	if x.hasOffset {
		panic(errors.Errorf("prevが存在しないNode(%v)のGetPrevEmptyNodeを取ろうとした", x))
	}
	return x.base
}

/////

func (x *Node) setBase(y Index) {
	log.Printf("nodes[%d](%s).%s %d -> %d", x.id, x, x.baseLabel(), x.base, y)
	x.base = y
}
func (x *Node) setCheck(y Index) {
	log.Printf("nodes[%d](%s).%s %d -> %d", x.id, x, x.checkLabel(), x.check, y)
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
