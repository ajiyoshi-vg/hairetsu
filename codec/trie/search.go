package trie

import (
	"iter"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Searcher[T any, DA doublearray.Nodes] struct {
	enc codec.Encoder[T]
	da  DA
}

func NewSearcher[T any, DA doublearray.Nodes](enc codec.Encoder[T], da DA) *Searcher[T, DA] {
	return &Searcher[T, DA]{enc: enc, da: da}
}

func (s *Searcher[T, DA]) ExactMatchSearch(x T) (node.Index, error) {
	return ExactMatchSearch(s.da, s.enc.Iter(x))
}

func (s *Searcher[T, DA]) CommonPrefixSearch(x T) ([]node.Index, error) {
	return CommonPrefixSearch(s.da, s.enc.Iter(x))
}

/*

// today, we can't write codes like this:
// refs. https://github.com/golang/go/issues/51053

type codes interface {
	iter.Seq2[int, word.Code] | []word.Code
}

func exactMatchSearch[T codes](cs T) {
	for i, c := range cs {
		fmt.Println(i, c)
	}
}
*/

func ExactMatchSearch[DA doublearray.Nodes](da DA, seq iter.Seq[word.Code]) (node.Index, error) {
	target, parent, err := doublearray.InitialTarget(da)
	if err != nil {
		return 0, err
	}
	for c := range seq {
		target, parent, err = doublearray.NextTarget(da, c, target, parent)
		if err != nil {
			return 0, err
		}
	}
	if !target.IsTerminal() {
		return 0, doublearray.ErrNotATerminal
	}
	data, err := da.At(target.GetChild(word.EOS))
	if err != nil {
		return 0, err
	}
	return data.GetOffset(), nil
}

func CommonPrefixSearch[DA doublearray.Nodes](da DA, seq iter.Seq[word.Code]) ([]node.Index, error) {
	var ret []node.Index
	target, parent, err := doublearray.InitialTarget(da)
	if err != nil {
		return nil, err
	}

	for c := range seq {
		target, parent, err = doublearray.NextTarget(da, c, target, parent)
		if err != nil {
			return nil, err
		}
		if target.IsTerminal() {
			data, err := da.At(target.GetChild(word.EOS))
			if err != nil {
				return ret, nil
			}
			ret = append(ret, data.GetOffset())
		}
	}
	return ret, nil
}
