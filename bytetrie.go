package hairetsu

import (
	"bufio"
	"io"
	"os"

	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/keytree"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type ByteTrie struct {
	data *da.DoubleArray
	b    da.Bytes
}

func (t *ByteTrie) ExactMatchSearch(key []byte) (node.Index, error) {
	return t.b.ExactMatchSearch(t.data, key)
}

func (t *ByteTrie) CommonPrefixSearch(key []byte) ([]node.Index, error) {
	return t.data.CommonPrefixSearch(word.FromBytes(key))
}

func (t *ByteTrie) WriteTo(w io.Writer) (int64, error) {
	return t.data.WriteTo(w)
}

func NewByteTrie(data *da.DoubleArray) *ByteTrie {
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

func (b *ByteTrieBuilder) Build(xs [][]byte) (*ByteTrie, error) {
	ret := da.New()
	ks, err := keytree.FromBytes(xs)
	if err != nil {
		return nil, err
	}
	if err := b.builder.Build(ret, ks); err != nil {
		return nil, err
	}
	return &ByteTrie{data: ret}, nil
}

func (b *ByteTrieBuilder) BuildFromFile(path string) (*ByteTrie, error) {
	ret := da.New()

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ks := keytree.New()

	r := bufio.NewScanner(file)
	for i := 0; r.Scan(); i++ {
		line := r.Text()
		key := word.FromBytes([]byte(line))
		if err := ks.Put(key, uint32(i)); err != nil {
			return nil, err
		}
	}
	if err := b.builder.Build(ret, ks); err != nil {
		return nil, err
	}
	return &ByteTrie{data: ret}, nil
}
