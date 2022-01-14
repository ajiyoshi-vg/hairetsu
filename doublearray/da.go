package doublearray

import (
	"fmt"
	"io"
	"strings"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/pkg/errors"
)

type DoubleArray struct {
	nodes []node.Node
}

func New() *DoubleArray {
	return &DoubleArray{
		nodes: make([]node.Node, 10),
	}
}

func FromArray(xs []uint64) *DoubleArray {
	nodes := make([]node.Node, len(xs))
	for i, x := range xs {
		nodes[i] = node.Node(x)
	}
	return &DoubleArray{
		nodes: nodes,
	}
}

func (da *DoubleArray) Array() []uint64 {
	ret := make([]uint64, len(da.nodes))
	for i, x := range da.nodes {
		ret[i] = uint64(x)
	}
	return ret
}

func (da *DoubleArray) ExactMatchSearch(cs word.Word) (node.Index, error) {
	var index node.Index
	length := node.Index(len(da.nodes))

	for _, c := range cs {
		next := da.at(index).GetOffset().Forward(c)
		if next >= length || !da.at(next).IsChildOf(index) {
			return 0, fmt.Errorf("ExactMatchSearch(%v) : error broken index", cs)
		}
		index = next
	}
	if !da.at(index).IsTerminal() {
		return 0, fmt.Errorf("ExactMatchSearch(%v) : not stored", cs)
	}
	data := da.at(index).GetOffset().Forward(word.EOS)
	if data >= length || !da.at(data).IsChildOf(index) {
		return 0, fmt.Errorf("ExactMatchSearch(%v) : error broken data node", cs)
	}
	return da.at(data).GetOffset(), nil
}

func (da *DoubleArray) CommonPrefixSearch(cs word.Word) ([]node.Index, error) {
	ret := make([]node.Index, 0, 10)
	var index node.Index
	length := node.Index(len(da.nodes))

	for _, c := range cs {
		next := da.at(index).GetOffset().Forward(c)
		if next >= length || !da.at(next).IsChildOf(index) {
			return ret, nil
		}
		index = next

		if da.at(index).IsTerminal() {
			data := da.at(index).GetOffset().Forward(word.EOS)
			if data >= length || !da.at(data).IsChildOf(index) {
				err := fmt.Errorf("CommonPrefixSearch(%v) : error broken data node", cs)
				return nil, err
			}
			ret = append(ret, da.at(data).GetOffset())
		}
	}
	return ret, nil
}

func (da *DoubleArray) Stat() Stat {
	return newStat(da)
}

func (da *DoubleArray) WriteTo(w io.Writer) (int64, error) {
	var ret int64
	for _, node := range da.nodes {
		buf, err := node.MarshalBinary()
		if err != nil {
			return ret, err
		}
		n, err := w.Write(buf)
		ret += int64(n)
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}

func (da *DoubleArray) Trace(cs word.Word) string {
	var index node.Index
	var err error
	ss := make([]string, 0, len(cs)+2)
	ss = append(ss, fmt.Sprintf("search:%v", cs))
	for _, c := range cs {
		ss = append(ss, da.debug(index, c))

		tmp := index
		index, err = da.traverse(index, c)
		if err != nil {
			ss = append(ss, fmt.Sprintf(
				"got error:%s at(%d) branch:%d",
				err,
				tmp,
				c,
			))
			return strings.Join(ss, "\n")
		}
	}
	ss = append(ss, fmt.Sprintf("%s", da.debug(index, word.EOS)))
	if da.at(index).IsTerminal() {
		data := da.at(index).GetOffset().Forward(word.EOS)
		if int(data) < len(da.nodes) {
			ss = append(ss, fmt.Sprintf("%s value", da.debug(data, word.EOS)))
		}
	}
	return strings.Join(ss, "\n")
}

func (da *DoubleArray) debug(i node.Index, c word.Code) string {
	return fmt.Sprintf(
		"at(%d):%s branch:%d forward:%d",
		i,
		da.at(i),
		c,
		da.at(i).GetOffset().Forward(c),
	)
}

func (da *DoubleArray) traverse(index node.Index, branch word.Code) (node.Index, error) {
	offset := da.at(index).GetOffset()
	next := offset.Forward(branch)
	if int(next) >= len(da.nodes) {
		return 0, errors.Errorf(
			"out of range nodes[%d] index:%d(%v) branch:%v",
			next,
			index,
			da.at(index),
			branch,
		)
	}
	if !da.at(next).IsChildOf(index) {
		return 0, errors.Errorf(
			"traverse fail node[%d](%v) is not child of node[%d](%v) branch:%d",
			next,
			da.at(next),
			index,
			da.at(index),
			branch,
		)
	}
	return next, nil
}

func (da *DoubleArray) getValue(term node.Index) (node.Index, error) {
	data, err := da.traverse(term, word.EOS)
	if err != nil {
		return 0, err
	}
	return da.at(data).GetOffset(), nil
}

func (da *DoubleArray) at(i node.Index) *node.Node {
	return &da.nodes[i]
}
