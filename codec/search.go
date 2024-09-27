package codec

import (
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Searcher[T any, DA doublearray.Nodes] struct {
	enc Encoder[T]
	da  DA
}

func NewSearcher[T any, DA doublearray.Nodes](enc Encoder[T], da DA) *Searcher[T, DA] {
	return &Searcher[T, DA]{enc: enc, da: da}
}

func (s *Searcher[T, DA]) ExactMatchSearch(x T) (node.Index, error) {
	target, parent, err := doublearray.InitialTarget(s.da)
	if err != nil {
		return 0, err
	}
	for c := range s.enc.Iter(x) {
		target, parent, err = doublearray.NextTarget(s.da, c, target, parent)
		if err != nil {
			return 0, err
		}
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

func (s *Searcher[T, DA]) CommonPrefixSearch(x T) ([]node.Index, error) {
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
