package hairetsu

import (
	"bytes"
	"io"

	dict "github.com/ajiyoshi-vg/hairetsu/bytes"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/node"
)

type DictTrie struct {
	data da.Nodes
	dict da.BytesDict
}

type DictTrieBuilder struct {
	builder *da.Builder
}

func NewDictTrie(data da.Nodes, dict dict.Dict) *DictTrie {
	return &DictTrie{
		data: data,
		dict: da.BytesDict(dict),
	}
}

func (t *DictTrie) ExactMatchSearch(key []byte) (node.Index, error) {
	return t.dict.ExactMatchSearch(t.data, key)
}

func (t *DictTrie) CommonPrefixSearch(key []byte) ([]node.Index, error) {
	return t.dict.CommonPrefixSearch(t.data, key)
}

func (t *DictTrie) WriteTo(w io.Writer) (int64, error) {
	ret, err := t.GetDict().WriteTo(w)
	if err != nil {
		return ret, err
	}
	n, err := t.data.WriteTo(w)
	ret += n
	return ret, err
}

func (t *DictTrie) ReadFrom(r io.Reader) (int64, error) {
	buf, err := io.ReadAll(r)
	ret := int64(len(buf))
	if err != nil {
		return ret, err
	}

	d := dict.Dict{}
	if err := d.UnmarshalBinary(buf[0:dict.Size]); err != nil {
		return ret, err
	}

	data := doublearray.New()
	rData := bytes.NewReader(buf[dict.Size:])

	if _, err := data.ReadFrom(rData); err != nil {
		return ret, err
	}

	t.dict = doublearray.BytesDict(d)
	t.data = data

	return ret, nil
}

func (t *DictTrie) GetDict() dict.Dict {
	return dict.Dict(t.dict)
}

func NewDictTrieBuilder(opt ...da.Option) *DictTrieBuilder {
	return &DictTrieBuilder{
		builder: da.NewBuilder(opt...),
	}
}

func (b *DictTrieBuilder) BuildFromSlice(xs [][]byte) (*DictTrie, error) {
	f := b.builder.Factory()
	dict, err := dict.FromSlice(xs, f)
	if err != nil {
		return nil, err
	}
	return buildDictTrie(f, dict)
}

func (b *DictTrieBuilder) BuildFromLines(r io.ReadSeeker) (*DictTrie, error) {
	f := b.builder.Factory()
	dict, err := dict.FromReadSeeker(r, f)
	if err != nil {
		return nil, err
	}
	return buildDictTrie(f, dict)
}

func buildDictTrie(f *da.Factory, dict dict.Dict) (*DictTrie, error) {
	trie, err := f.Done()
	if err != nil {
		return nil, err
	}
	return NewDictTrie(trie, dict), nil
}
