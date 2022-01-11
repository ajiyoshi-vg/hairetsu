package hairetsu

import (
	"fmt"
	"sort"

	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/keyset"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type RuneDict map[rune]word.Code

type RuneTrie struct {
	data *da.DoubleArray
	dict RuneDict
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
	ret := da.New(len(xs) * 2)
	dict := newRuneDict(xs)
	ks, err := b.keyset(xs, dict)
	if err != nil {
		return nil, err
	}
	if err := b.builder.Build(ret, ks); err != nil {
		return nil, err
	}
	return &RuneTrie{data: ret, dict: dict}, nil
}

func (*RuneTrieBuilder) keyset(ss []string, d RuneDict) (keyset.KeySet, error) {
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

func newRuneDict(ss []string) RuneDict {
	runeCount := make(map[rune]uint32, len(ss))
	for _, s := range ss {
		for _, r := range s {
			runeCount[r] += 1
		}
	}

	type tmp struct {
		r rune
		n uint32
	}

	buf := make([]tmp, 0, len(runeCount))
	for r, n := range runeCount {
		buf = append(buf, tmp{r: r, n: n})
	}

	sort.Slice(buf, func(i, j int) bool {
		return buf[i].n > buf[j].n
	})

	ret := make(RuneDict, len(buf))
	for i, x := range buf {
		ret[x.r] = word.Code(i + 1)
	}
	return ret
}

func (d RuneDict) Word(s string) word.Word {
	ret := make(word.Word, 0, len(s))
	for _, r := range s {
		c, ok := d[r]
		if !ok {
			ret = append(ret, word.NONE)
		} else {
			ret = append(ret, c)
		}
	}
	return ret
}

func (d RuneDict) StrictWord(s string) (word.Word, error) {
	ret := make(word.Word, 0, len(s))
	for _, r := range s {
		c, ok := d[r]
		if !ok {
			return nil, fmt.Errorf("unknown rune(%c)", r)
		} else {
			ret = append(ret, c)
		}
	}
	return ret, nil
}
