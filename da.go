package hairetsu

import (
	"strings"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/pkg/errors"
)

type DoubleArray struct {
	nodes   []node.Node
	factory nodeFactory
}

func (da *DoubleArray) init(after int) {
	if after == 0 {
		da.nodes[0] = da.factory.root()
		after = 1
	}

	for i := after; i < len(da.nodes); i++ {
		da.nodes[i] = da.factory.node(i)
	}
}

func (da *DoubleArray) ExactMatchSearch(xs []byte) (node.Index, error) {
	cs := word.FromBytes(xs)
	return da.lookup(cs)
}
func (da *DoubleArray) CommonPrefixSearch(xs []byte) ([]node.Index, error) {
	cs := word.FromBytes(xs)
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

func (da *DoubleArray) extend() {
	max := len(da.nodes)
	da.nodes = append(da.nodes, make([]node.Node, len(da.nodes))...)
	da.init(max)
}

func (da *DoubleArray) ensure(i node.Index) {
	for len(da.nodes) <= int(i) {
		da.extend()
	}
}

func (da *DoubleArray) popNode(i node.Index) {
	// これから nodes[i] を使うための準備
	// nodes[i] を prev/next にしているnodeから node[i]を取り除く

	prev := da.nodes[i].GetPrevEmptyNode()
	next := da.nodes[i].GetNextEmptyNode()

	// next にアクセスできるように、必要があれば拡張
	da.ensure(next)

	// 1. nodes[i].prev の next に nodes[i].next を繋ぐ
	da.nodes[prev].SetNextEmptyNode(next)
	// 2. nodes[i].next の prev に nodes[i].prev を繋ぐ
	da.nodes[next].SetPrevEmptyNode(prev)
}

func (da *DoubleArray) traverse(index node.Index, branch word.Code) (node.Index, error) {
	offset := da.nodes[index].GetOffset()
	next := offset.Forward(branch)
	if int(next) >= len(da.nodes) || !da.nodes[next].IsChildOf(index) {
		return 0, errors.WithMessagef(da.nodeError(), "branch:%v", branch)
	}
	return next, nil
}

//FIXME
func (da *DoubleArray) nodeError() error {
	ss := make([]string, 0, len(da.nodes))
	for _, node := range da.nodes {
		ss = append(ss, node.String())
	}
	return errors.New(strings.Join(ss, "\n"))
}

func (da *DoubleArray) findValidOffset(cs word.Word) node.Index {
	index, offset := da.findOffset(da.nodes[0].GetNextEmptyNode(), cs[0])

	// offset からcs を全部格納可能なところを探す
	for i := 0; i < len(cs); i++ {
		next := offset.Forward(cs[i])

		if int(next) >= len(da.nodes) {
			break
		}

		if int(index) >= len(da.nodes) {
			break
		}

		if da.nodes[next].HasParent() {
			// 使用済みだった
			// 次の未使用ノードを試す
			index, offset = da.findOffset(da.nodes[index].GetNextEmptyNode(), cs[0])
			// cs[0] からやりなおし
			i = 0
		}
	}
	return offset
}
func (da *DoubleArray) findOffset(index node.Index, branch word.Code) (node.Index, node.Index) {
	for {
		offset, err := index.Backward(branch)
		if err == nil {
			return index, offset
		}
		da.ensure(index)
		index = da.nodes[index].GetNextEmptyNode()
	}
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

func (da *DoubleArray) lookup(cs word.Word) (node.Index, error) {
	index, err := da.getIndex(cs)
	if err != nil {
		return 0, err
	}
	if da.nodes[index].IsTerminal() {
		return da.getValue(index)
	}
	return 0, errors.WithMessagef(da.nodeError(), "not terminal. lookup:%v index:%d", cs, index)
}

func (da *DoubleArray) getValue(term node.Index) (node.Index, error) {
	data, err := da.traverse(term, word.EOS)
	if err != nil {
		return 0, err
	}
	return da.nodes[data].GetOffset(), nil
}
