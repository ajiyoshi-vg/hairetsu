package node

import (
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/pkg/errors"
)

type Index uint32

type NodeInterface interface {
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

	String() string
}

func (x Index) Forward(c word.Code) Index {
	return x + Index(c)
}

// Backword - offset.Forward(c) == x となるようなoffsetを返す
func (x Index) Backward(c word.Code) (Index, error) {
	if x < Index(c) {
		return 0, errors.Errorf("can't backword from %d by %d", x, c)
	}
	offset := x - Index(c)
	return offset, nil
}
