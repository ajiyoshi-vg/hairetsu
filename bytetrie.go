package hairetsu

import (
	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/keytree"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type ByteTrie struct {
	data *da.DoubleArray
}

func NewByteTrie(data *da.DoubleArray) *ByteTrie {
	return &ByteTrie{data: data}
}

type ByteTrieBuilder struct {
	builder *da.Builder
}

func NewByteTrieBuilder() *ByteTrieBuilder {
	return &ByteTrieBuilder{
		builder: da.NewBuilder(),
	}
}

func (b *ByteTrieBuilder) Build(xs [][]byte) (*ByteTrie, error) {
	ret := da.New(len(xs) * 2)
	ks := keytree.FromBytes(xs)
	if err := b.builder.Build(ret, ks); err != nil {
		return nil, err
	}
	return &ByteTrie{data: ret}, nil
}

func (t *ByteTrie) ExactMatchSearch(key []byte) (node.Index, error) {
	return t.data.ExactMatchSearch(word.FromBytes(key))
}

func (t *ByteTrie) CommonPrefixSearch(key []byte) ([]node.Index, error) {
	return t.data.CommonPrefixSearch(word.FromBytes(key))
}
