package keyset

import (
	"log"
	"sort"

	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Item struct {
	Key word.Word
	Val uint32
}

type KeySet []Item

type Callback func(word.Word, []word.Code, []uint32) error

func FromBytes(xs [][]byte) KeySet {
	ret := make(KeySet, 0, len(xs))
	for i, x := range xs {
		ret = append(ret, Item{
			Key: word.FromBytes(x),
			Val: uint32(i),
		})
	}
	return ret
}

func New(data []word.Word) KeySet {
	ret := make([]Item, 0, len(data))
	for i, x := range data {
		ret = append(ret, Item{Key: x, Val: uint32(i)})
	}
	return ret
}

func (ks KeySet) less(i, j int) bool {
	return word.Compare(ks[i].Key, ks[j].Key) < 0
}

func (ks KeySet) Sort() {
	if !sort.SliceIsSorted(ks, ks.less) {
		sort.Slice(ks, ks.less)
	}
}

func (ks KeySet) LeafNum() int {
	return len(ks)
}

func (ks KeySet) Walk(f func(word.Word, []word.Code, []uint32) error) error {
	ks.Sort()
	log.Println("sorted")
	return ks.walkTrieNode(0, len(ks), 0, f)
}

func (ks KeySet) walkTrieNode(begin, end, depth int, f func(word.Word, []word.Code, []uint32) error) error {
	// apply callback first
	if err := f(ks.trieNode(begin, end, depth)); err != nil {
		return err
	}

	for begin < end {
		if ks[begin].Key.At(depth) != word.EOS {
			break
		}
		begin++
	}
	if begin == end {
		return nil
	}

	lastBegin := begin
	lastLabel := ks[begin].Key[depth]
	begin++
	for begin < end {
		label := ks[begin].Key[depth]
		if label != lastLabel {
			if err := ks.walkTrieNode(lastBegin, begin, depth+1, f); err != nil {
				return err
			}
			lastBegin = begin
			lastLabel = label
		}
		begin++
	}
	return ks.walkTrieNode(lastBegin, end, depth+1, f)
}

func (ks KeySet) trieNode(begin, end, depth int) (word.Word, []word.Code, []uint32) {
	prefix := ks[begin].Key[0:depth]
	branch := make([]word.Code, 0, end)
	values := make([]uint32, 0, end)
	for i := begin; i < end; i++ {
		branch = append(branch, ks[i].Key.At(depth))
		values = append(values, ks[i].Val)
	}
	return prefix, branch, values
}
