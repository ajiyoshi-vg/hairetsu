package doublearray

import (
	"strings"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/pkg/errors"
)

type DoubleArray struct {
	nodes []node.Node
}

func New(initial int) *DoubleArray {
	return &DoubleArray{
		nodes: make([]node.Node, initial),
	}
}

func (da *DoubleArray) ExactMatchSearch(cs word.Word) (node.Index, error) {
	index, err := da.getIndex(cs)
	if err != nil {
		return 0, err
	}
	if da.nodes[index].IsTerminal() {
		return da.getValue(index)
	}
	return 0, errors.WithMessagef(da.nodeError(), "not terminal. lookup:%v index:%d", cs, index)
}

func (da *DoubleArray) CommonPrefixSearch(cs word.Word) ([]node.Index, error) {
	ret := make([]node.Index, 0, 10)

	var index node.Index
	var err error
	for _, c := range cs {
		index, err = da.traverse(index, c)
		if err != nil {
			return ret, nil
		}

		if da.nodes[index].IsTerminal() {
			val, err := da.getValue(index)
			if err != nil {
				return nil, err
			}
			ret = append(ret, val)
		}
	}
	return ret, nil
}

func (da *DoubleArray) Stat() Stat {
	return newStat(da)
}

func (da *DoubleArray) traverse(index node.Index, branch word.Code) (node.Index, error) {
	offset := da.nodes[index].GetOffset()
	next := offset.Forward(branch)
	if int(next) >= len(da.nodes) || !da.nodes[next].IsChildOf(index) {
		return 0, errors.WithMessagef(da.nodeError(), "branch:%v", branch)
	}
	return next, nil
}

func (da *DoubleArray) getValue(term node.Index) (node.Index, error) {
	data, err := da.traverse(term, word.EOS)
	if err != nil {
		return 0, err
	}
	return da.nodes[data].GetOffset(), nil
}

//FIXME
func (da *DoubleArray) nodeError() error {
	ss := make([]string, 0, len(da.nodes))
	for _, node := range da.nodes {
		ss = append(ss, node.String())
	}
	return errors.New(strings.Join(ss, "\n"))
}

func (da *DoubleArray) getIndex(cs word.Word) (node.Index, error) {
	var index node.Index
	var err error
	for _, c := range cs {
		index, err = da.traverse(index, c)
		if err != nil {
			return 0, errors.WithMessagef(err, "word:%v", cs)
		}
	}
	return index, nil
}
