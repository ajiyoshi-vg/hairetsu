package hairetsu

import (
	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/keyset"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/runedict"
)

type RuneTrie struct {
	data *da.DoubleArray
	dict runedict.RuneDict
}

type RuneTrieBuilder struct {
	builder *da.Builder
}

func (t *RuneTrie) ExactMatchSearch(key string) (node.Index, error) {
	return t.data.ExactMatchSearch(t.dict.Word(key))
}

func (t *RuneTrie) CommonPrefixSearch(key string) ([]node.Index, error) {
	return t.data.CommonPrefixSearch(t.dict.Word(key))
}

func NewRuneTrieBuilder() *RuneTrieBuilder {
	return &RuneTrieBuilder{
		builder: da.NewBuilder(),
	}
}

func (b *RuneTrieBuilder) Build(xs []string) (*RuneTrie, error) {
	ret := da.New()
	dict := runedict.New(xs)
	ks, err := b.keyset(xs, dict)
	if err != nil {
		return nil, err
	}
	if err := b.builder.Build(ret, ks); err != nil {
		return nil, err
	}
	return &RuneTrie{data: ret, dict: dict}, nil
}

func (*RuneTrieBuilder) keyset(ss []string, d runedict.RuneDict) (keyset.KeySet, error) {
	ret := make(keyset.KeySet, 0, len(ss))
	for i, s := range ss {
		w, err := d.StrictWord(s)
		if err != nil {
			return nil, err
		}
		ret = append(ret, keyset.Item{Key: w, Val: uint32(i)})
	}
	return ret, nil
}
