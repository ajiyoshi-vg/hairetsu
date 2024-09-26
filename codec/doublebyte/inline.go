package doublebyte

import (
	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type InlineSearcher[T Dict] struct {
	da   da.Nodes
	dict T
}

func NewInlineSearcher[T Dict](da da.Nodes, dict T) *InlineSearcher[T] {
	return &InlineSearcher[T]{da: da, dict: dict}
}

func (x *InlineSearcher[T]) ExactMatchSearch(key []byte) (node.Index, error) {
	target, parent, err := da.InitialTarget(x.da)
	if err != nil {
		return 0, err
	}
	n := len(key)
	for i := 0; i < n; i += 2 {
		val := uint16(key[i])
		if i+1 < n {
			val |= uint16(key[i+1]) << 8
		}
		code := x.dict.Code(val)
		target, parent, err = da.NextTarget(x.da, code, target, parent)
		if err != nil {
			return 0, err
		}
	}
	if len(key)%2 == 1 {
		target, _, err = da.NextTarget(x.da, word.Backspace, target, parent)
		if err != nil {
			return 0, err
		}
	}
	if !target.IsTerminal() {
		return 0, err
	}
	data, err := x.da.At(target.GetChild(word.EOS))
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}
