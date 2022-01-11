package u64

import (
	"fmt"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/pkg/errors"
)

/*

value bit layout
| 00 01 02 03 04 ...   32 33 34 ....     63 |
  ^     ^  ^              ^  ^
  |     |  |              |  |
  |     |  base(30bit)    |  check(30bit)
  |     hassOffset        hasParent
  | isTerminal

*/

/*o

| 00 01 02 03 04 ...   32 33 34 ....     63 |
                           0  1  1 ..     1
                        0  1  1 ..     1  0  << 1
         0  1  1 ...                         << 31
*/

const (
	isTerminal = 1 << 63
	hasOffset  = 1 << 61
	hasParent  = 1 << 30

	baseMask  = ((uint64(1) << 29) - 1) << 31
	checkMask = uint64(1)<<29 - 1
)

type Node uint64

func Root() Node {
	var ret Node
	ret.SetNextEmptyNode(1)
	return ret
}

func New(i int) Node {
	//XX
	var ret Node
	ret.SetPrevEmptyNode(node.Index(i - 1))
	ret.SetNextEmptyNode(node.Index(i + 1))
	return ret
}

func (x Node) GetOffset() node.Index {
	return x.getBase()
}

func (x *Node) SetOffset(i node.Index) {
	val := ^baseMask&uint64(*x) | uint64(i)<<31 | hasOffset
	*x = Node(val)
}

func (x *Node) Terminate() {
	*x |= isTerminal
}
func (x Node) IsTerminal() bool {
	return x&isTerminal > 0
}

func (x *Node) SetParent(i node.Index) {
	val := ^checkMask&uint64(*x) | uint64(i) | hasParent
	*x = Node(val)
}
func (x Node) GetParent() node.Index {
	return x.getCheck()
}
func (x Node) IsChildOf(parent node.Index) bool {
	if !x.HasParent() {
		return false
	}
	return x.GetParent() == parent
}
func (x Node) HasParent() bool {
	return x&hasParent > 0
}
func (x Node) GetNextEmptyNode() node.Index {
	if x.HasParent() {
		panic(errors.Errorf("emptyでないNode(%v)のnextEmptyNodeを取ろうとした", x))
	}
	return x.getCheck()
}
func (x Node) GetPrevEmptyNode() node.Index {
	if x&hasOffset > 0 {
		panic(errors.Errorf("prevが存在しないNode(%v)のGetPrevEmptyNodeを取ろうとした", x))
	}
	return x.getBase()
}

func (x *Node) SetNextEmptyNode(i node.Index) {
	val := ^checkMask&uint64(*x) | uint64(i)
	*x = Node(val)
}

func (x *Node) SetPrevEmptyNode(i node.Index) {
	val := ^baseMask&uint64(*x) | uint64(i)<<31
	*x = Node(val)
}

func (x Node) getBase() node.Index {
	ret := (uint64(x) & baseMask) >> 31
	return node.Index(ret)
}
func (x Node) getCheck() node.Index {
	ret := (uint64(x) & checkMask)
	return node.Index(ret)
}
func (x Node) baseLabel() string {
	if x&hasOffset > 0 {
		return "base"
	}
	return "prev"
}
func (x Node) checkLabel() string {
	if x&hasParent > 0 {
		return "check"
	}
	return "next"
}
func (x Node) String() string {
	ret := fmt.Sprintf("{%s:%d, %s:%d}",
		x.baseLabel(),
		x.getBase(),
		x.checkLabel(),
		x.getCheck(),
	)
	if x.IsTerminal() {
		ret += "#"
	}
	return ret
}
