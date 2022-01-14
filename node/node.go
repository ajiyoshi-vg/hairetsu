// +build !verbose

package node

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"math/rand"
	"reflect"

	"github.com/pkg/errors"
)

/*

Node bit layout (uint64)
| 00 01 02 03 04 ...    32 33 34 ....     63 |
   ^     ^  ^            ^  ^  ^           ^
   |     |  |  (30bit)   |  |  |  (30bit)  |
   |     |  |_base_______|  |  |_check_____|
   |     |                  |
   |     |_hasOffset        |_hasParent
   |
   |_isTerminal

*/

const (
	isTerminal = 1 << 63
	hasOffset  = 1 << 61
	hasParent  = 1 << 30

	baseMask  = ((uint64(1) << 29) - 1) << 31
	checkMask = uint64(1)<<29 - 1
)

type Interface interface {
	GetOffset() Index
	SetOffset(Index)

	Terminate()
	IsTerminal() bool

	SetParent(Index)
	IsChildOf(Index) bool

	IsUsed() bool

	GetNextEmptyNode() (Index, error)
	GetPrevEmptyNode() (Index, error)
	SetNextEmptyNode(Index) error
	SetPrevEmptyNode(Index) error

	Reset(int)
	String() string
}

type Node uint64

var (
	_ Interface                = (*Node)(nil)
	_ encoding.BinaryMarshaler = (*Node)(nil)
)

func Root() Node {
	var ret Node
	ret.SetNextEmptyNode(1)
	return ret
}

func New(i int) Node {
	//XX
	var ret Node
	ret.SetPrevEmptyNode(Index(i - 1))
	ret.SetNextEmptyNode(Index(i + 1))
	return ret
}

func (x *Node) Reset(i int) {
	if i == 0 {
		*x = Root()
	} else {
		*x = New(i)
	}
}

func (x Node) GetOffset() Index {
	return x.getBase()
}

func (x *Node) SetOffset(i Index) {
	val := ^baseMask&uint64(*x) | uint64(i)<<31 | hasOffset
	*x = Node(val)
}

func (x *Node) Terminate() {
	*x |= isTerminal
}
func (x Node) IsTerminal() bool {
	return x&isTerminal > 0
}

func (x *Node) SetParent(i Index) {
	val := ^checkMask&uint64(*x) | uint64(i) | hasParent
	*x = Node(val)
}
func (x Node) GetParent() Index {
	return x.getCheck()
}
func (x Node) IsChildOf(parent Index) bool {
	if !x.HasParent() {
		return false
	}
	return x.GetParent() == parent
}
func (x Node) IsUsed() bool {
	return x.HasOffset() || x.HasParent()
}
func (x Node) HasParent() bool {
	return x&hasParent > 0
}
func (x Node) HasOffset() bool {
	return x&hasOffset > 0
}
func (x Node) GetNextEmptyNode() (Index, error) {
	if x.HasParent() {
		return 0, errors.Errorf("try to GetNextEmptyNode of used Node(%s)", x)
	}
	return x.getCheck(), nil
}
func (x Node) GetPrevEmptyNode() (Index, error) {
	if x.HasOffset() {
		return 0, errors.Errorf("try to GetPrevEmptyNode of used Node(%s)", x)
	}
	return x.getBase(), nil
}

func (x *Node) SetNextEmptyNode(i Index) error {
	if x.HasParent() {
		return errors.Errorf("try to SetNextEmptyNode of used Node(%s)", x)
	}
	val := ^checkMask&uint64(*x) | uint64(i)
	*x = Node(val)
	return nil
}

func (x *Node) SetPrevEmptyNode(i Index) error {
	if x.HasOffset() {
		return errors.Errorf("try to SetPrevEmptyNode of used Node(%s)", x)
	}
	val := ^baseMask&uint64(*x) | uint64(i)<<31
	*x = Node(val)
	return nil
}

func (x Node) getBase() Index {
	ret := (uint64(x) & baseMask) >> 31
	return Index(ret)
}
func (x Node) getCheck() Index {
	ret := (uint64(x) & checkMask)
	return Index(ret)
}
func (x Node) baseLabel() string {
	if x.HasOffset() {
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
func (x Node) MarshalBinary() ([]byte, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(x))
	return buf, nil
}
func (x *Node) UnmarshalBinary(data []byte) error {
	v := binary.BigEndian.Uint64(data)
	*x = Node(v)
	return nil
}
func (x Node) Generate(r *rand.Rand, size int) reflect.Value {
	ret := Node(r.Uint64())
	return reflect.ValueOf(ret)
}
