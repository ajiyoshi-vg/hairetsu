package hairetsu

import (
	"io"
	"iter"
	"slices"

	"github.com/ajiyoshi-vg/external/scan"
	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type ByteTrie struct {
	data da.Nodes
	b    da.Bytes
}

func (t *ByteTrie) ExactMatchSearch(key []byte) (node.Index, error) {
	return t.b.ExactMatchSearch(t.data, key)
}

func (t *ByteTrie) CommonPrefixSearch(key []byte) ([]node.Index, error) {
	return t.b.CommonPrefixSearch(t.data, key)
}

func (t *ByteTrie) WriteTo(w io.Writer) (int64, error) {
	return t.data.WriteTo(w)
}

func NewByteTrie(data da.Nodes) *ByteTrie {
	return &ByteTrie{data: data}
}

type ByteTrieBuilder struct {
	builder *da.Builder
}

func NewByteTrieBuilder(opt ...da.Option) *ByteTrieBuilder {
	return &ByteTrieBuilder{
		builder: da.NewBuilder(opt...),
	}
}

func (b *ByteTrieBuilder) BuildSlice(xs [][]byte) (*ByteTrie, error) {
	return b.StreamBuild(slices.Values(xs))
}

func (b *ByteTrieBuilder) StreamBuild(seq iter.Seq[[]byte]) (*ByteTrie, error) {
	x, err := b.builder.StreamBuild(item.FromByteSeq(seq))
	if err != nil {
		return nil, err
	}
	return NewByteTrie(x), nil
}

func (b *ByteTrieBuilder) BuildFromLines(r io.Reader) (*ByteTrie, error) {
	f := b.builder.Factory()
	var i uint32
	for x := range scan.ByteLines(r) {
		f.Put(item.New(word.FromBytes(x), i))
		i++
	}
	ret, err := f.Done()
	if err != nil {
		return nil, err
	}
	return NewByteTrie(ret), nil
}
