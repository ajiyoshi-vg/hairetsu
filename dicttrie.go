package hairetsu

import (
	bt "bytes"
	"io"
	"io/ioutil"
	"math"

	"github.com/ajiyoshi-vg/hairetsu/bytes"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/token"
)

type DictTrie struct {
	data da.Nodes
	dict da.BytesDict
}

type DictTrieBuilder struct {
	builder *da.Builder
}

func NewDictTrie(data da.Nodes, dict bytes.Dict) *DictTrie {
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
	return ret, nil
}

func (t *DictTrie) ReadFrom(r io.Reader) (int64, error) {
	buf, err := ioutil.ReadAll(r)
	ret := int64(len(buf))
	if err != nil {
		return ret, err
	}

	dict := bytes.Dict{}
	if err := dict.UnmarshalBinary(buf[0:math.MaxUint8]); err != nil {
		return ret, err
	}

	data := doublearray.New()
	rData := bt.NewReader(buf[math.MaxUint8:])

	if _, err := data.ReadFrom(rData); err != nil {
		return ret, err
	}

	t.dict = doublearray.BytesDict(dict)
	t.data = data

	return ret, nil
}

func (t *DictTrie) GetDict() bytes.Dict {
	return bytes.Dict(t.dict)
}

func NewDictTrieBuilder(opt ...da.Option) *DictTrieBuilder {
	return &DictTrieBuilder{
		builder: da.NewBuilder(opt...),
	}
}

func (b *DictTrieBuilder) Build(ks doublearray.Walker, dict bytes.Dict) (*DictTrie, error) {
	data := da.New()
	if err := b.builder.Build(data, ks); err != nil {
		return nil, err
	}
	return NewDictTrie(data, dict), nil
}

func (b *DictTrieBuilder) BuildFromLines(r io.Reader) (*DictTrie, error) {
	ks, dict, err := bytes.FromWalker(token.NewLinedBytes(r))
	if err != nil {
		return nil, err
	}
	return b.Build(ks, dict)
}
