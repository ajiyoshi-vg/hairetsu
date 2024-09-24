package hairetsu

import (
	"io"
	"iter"
	"slices"

	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/ajiyoshi-vg/hairetsu/node"
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
	x, err := da.StreamBuild(item.FromByteSeq(seq))
	if err != nil {
		return nil, err
	}
	return NewByteTrie(x), nil
}
