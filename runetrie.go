package hairetsu

import (
	"bufio"
	"io"
	"os"

	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/keyset"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/runedict"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type RuneTrie struct {
	data *da.DoubleArray
	dict Dict
}

type Dict interface {
	Word(string) word.Word
	StrictWord(string) (word.Word, error)
	MarshalText() (string, error)
}

type RuneTrieBuilder struct {
	builder *da.Builder
}

func NewRuneTrie(data *da.DoubleArray, dict runedict.RuneDict) *RuneTrie {
	return &RuneTrie{
		data: data,
		dict: dict,
	}
}

func (t *RuneTrie) ExactMatchSearch(key string) (node.Index, error) {
	return t.data.ExactMatchSearch(t.dict.Word(key))
}

func (t *RuneTrie) CommonPrefixSearch(key string) ([]node.Index, error) {
	return t.data.CommonPrefixSearch(t.dict.Word(key))
}

func (t *RuneTrie) WriteTo(w io.Writer) (int64, error) {
	return t.data.WriteTo(w)
}

func (t *RuneTrie) GetDict() Dict {
	return t.dict
}

func NewRuneTrieBuilder(opt ...da.Option) *RuneTrieBuilder {
	return &RuneTrieBuilder{
		builder: da.NewBuilder(opt...),
	}
}

func (b *RuneTrieBuilder) Build(xs []string) (*RuneTrie, error) {
	data := da.New()
	dict := runedict.New(xs)
	ks, err := b.keyset(xs, dict)
	if err != nil {
		return nil, err
	}
	if err := b.builder.Build(data, ks); err != nil {
		return nil, err
	}
	return &RuneTrie{data: data, dict: dict}, nil
}

func (b *RuneTrieBuilder) BuildFromFile(path string) (*RuneTrie, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := bufio.NewScanner(file)
	ss := make([]string, 0, 100)
	for r.Scan() {
		line := r.Text()
		ss = append(ss, line)
	}

	return b.Build(ss)
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
