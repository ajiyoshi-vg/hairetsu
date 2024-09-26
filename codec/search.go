package codec

import (
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Searcher[T any] struct {
	enc Encoder[T]
	da  doublearray.Nodes
}

func NewSearcher[T any](enc Encoder[T], da doublearray.Nodes) *Searcher[T] {
	return &Searcher[T]{enc: enc, da: da}
}

func (s *Searcher[T]) ExactMatchSearch(x T) (node.Index, error) {
	var parent node.Index
	target, err := s.da.At(parent)
	if err != nil {
		return 0, err
	}
	for c := range s.enc.Iter(x) {
		child := target.GetChild(c)
		target, err = s.da.At(child)
		if err != nil {
			return 0, err
		}
		if !target.IsChildOf(parent) {
			return 0, doublearray.ErrNotAChild
		}
		parent = child
	}
	if !target.IsTerminal() {
		return 0, doublearray.ErrNotATerminal
	}
	data, err := s.da.At(target.GetChild(word.EOS))
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}

func (s *Searcher[T]) CommonPrefixSearch(x T) ([]node.Index, error) {
	var ret []node.Index
	var parent node.Index
	target, err := s.da.At(parent)
	if err != nil {
		return ret, nil
	}

	for c := range s.enc.Iter(x) {
		child := target.GetChild(c)
		target, err = s.da.At(child)
		if err != nil {
			return ret, nil
		}
		if !target.IsChildOf(parent) {
			return ret, nil
		}
		parent = child
		if target.IsTerminal() {
			data, err := s.da.At(target.GetChild(word.EOS))
			if err != nil {
				return ret, nil
			}
			ret = append(ret, data.GetOffset())
		}
	}
	return ret, nil
}
