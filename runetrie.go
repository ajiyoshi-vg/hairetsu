package hairetsu

import (
	"io"

	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/keyset"
	"github.com/ajiyoshi-vg/hairetsu/keytree"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/runedict"
)

type RuneTrie struct {
	data da.Nodes
	dict da.Runes
}

type RuneTrieBuilder struct {
	builder *da.Builder
}

func NewRuneTrie(data da.Nodes, dict runedict.RuneDict) *RuneTrie {
	return &RuneTrie{
		data: data,
		dict: da.Runes(dict),
	}
}

func (t *RuneTrie) ExactMatchSearch(key string) (node.Index, error) {
	return t.dict.ExactMatchSearch(t.data, key)
}

func (t *RuneTrie) CommonPrefixSearch(key string) ([]node.Index, error) {
	return t.dict.CommonPrefixSearch(t.data, key)
}

func (t *RuneTrie) WriteTo(w io.Writer) (int64, error) {
	return t.data.WriteTo(w)
}

func (t *RuneTrie) GetDict() runedict.RuneDict {
	return runedict.RuneDict(t.dict)
}

func NewRuneTrieBuilder(opt ...da.Option) *RuneTrieBuilder {
	return &RuneTrieBuilder{
		builder: da.NewBuilder(opt...),
	}
}

func (b *RuneTrieBuilder) Build(ks doublearray.Walker, dict runedict.RuneDict) (*RuneTrie, error) {
	data := da.New()
	if err := b.builder.Build(data, ks); err != nil {
		return nil, err
	}
	return NewRuneTrie(data, dict), nil
}

func (b *RuneTrieBuilder) BuildSlice(xs []string) (*RuneTrie, error) {
	dict := runedict.New(xs)
	ks, err := b.keyset(xs, dict)
	if err != nil {
		return nil, err
	}
	return b.Build(ks, dict)
}

func (b *RuneTrieBuilder) BuildFromLines(r io.Reader) (*RuneTrie, error) {
	ks, dict, err := keytree.FromStringLines(r)
	if err != nil {
		return nil, err
	}
	return b.Build(ks, dict)
}

func (*RuneTrieBuilder) keyset(ss []string, d runedict.RuneDict) (keyset.KeySet, error) {
	ret := make(keyset.KeySet, 0, len(ss))
	for i, s := range ss {
		w, err := d.Word(s)
		if err != nil {
			return nil, err
		}
		ret = append(ret, keyset.Item{Key: w, Val: uint32(i)})
	}
	return ret, nil
}
