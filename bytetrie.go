package hairetsu

import (
	"io"

	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/keytree"
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
	ks, err := keytree.FromBytes(xs)
	if err != nil {
		return nil, err
	}
	return b.Build(ks)
}

func (b *ByteTrieBuilder) Build(ks da.Walker) (*ByteTrie, error) {
	x := da.New()
	if err := b.builder.Build(x, ks); err != nil {
		return nil, err
	}
	return NewByteTrie(x), nil
}
