package doublearray

import (
	"encoding/binary"
	"fmt"

	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"golang.org/x/exp/mmap"
)

type Mmap struct {
	r *mmap.ReaderAt
}

const nodeSize = 8

func NewMmap(path string) (*Mmap, error) {
	// will be closed by runtime.SetFinalizer
	r, err := mmap.Open(path)
	if err != nil {
		return nil, err
	}
	return &Mmap{r: r}, nil
}

func (da *Mmap) at(i node.Index) node.Node {
	s := make([]byte, nodeSize)
	n, err := da.r.ReadAt(s, int64(i)*nodeSize)
	if n != nodeSize {
		panic("bad size")
	}
	if err != nil {
		panic(err)
	}

	v := binary.BigEndian.Uint64(s)
	ret := node.Node(v)
	return ret
}

func (da *Mmap) length() int {
	return da.r.Len()
}

func (da *Mmap) ExactMatchSearch(cs word.Word) (node.Index, error) {
	var index node.Index
	length := node.Index(da.length())

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

func (da *Mmap) CommonPrefixSearch(cs word.Word) ([]node.Index, error) {
	var ret []node.Index
	var index node.Index
	length := node.Index(da.length())

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
