package hairetsu

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"

	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/keyset"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/runes"
	"github.com/ajiyoshi-vg/hairetsu/token"
)

type RuneTrie struct {
	data da.Nodes
	dict da.RunesDict
}

type RuneTrieBuilder struct {
	builder *da.Builder
}

func NewRuneTrie(data da.Nodes, dict runes.Dict) *RuneTrie {
	return &RuneTrie{
		data: data,
		dict: da.RunesDict(dict),
	}
}

func (t *RuneTrie) ExactMatchSearch(key string) (node.Index, error) {
	return t.dict.ExactMatchSearch(t.data, key)
}

func (t *RuneTrie) CommonPrefixSearch(key string) ([]node.Index, error) {
	return t.dict.CommonPrefixSearch(t.data, key)
}

const runeTrieHeader = 4

func (t *RuneTrie) WriteTo(w io.Writer) (int64, error) {
	dict := &bytes.Buffer{}
	n, err := t.GetDict().WriteTo(dict)
	if err != nil {
		return 0, err
	}

	buf := make([]byte, runeTrieHeader)
	binary.BigEndian.PutUint32(buf, uint32(n))
	size := bytes.NewBuffer(buf)

	data := &bytes.Buffer{}
	_, err = t.data.WriteTo(data)

	return io.Copy(w, io.MultiReader(size, dict, data))
}

func (t *RuneTrie) ReadFrom(r io.Reader) (int64, error) {
	buf, err := ioutil.ReadAll(r)
	ret := int64(len(buf))
	if err != nil {
		return ret, err
	}
	if ret < runeTrieHeader {
		return ret, io.ErrUnexpectedEOF
	}
	size := binary.BigEndian.Uint32(buf[0:runeTrieHeader])

	dict := runes.Dict{}
	rDict := bytes.NewReader(buf[runeTrieHeader : runeTrieHeader+size])
	if _, err := dict.ReadFrom(rDict); err != nil {
		return ret, err
	}
	data := doublearray.New()
	rData := bytes.NewReader(buf[runeTrieHeader+size:])

	if _, err := data.ReadFrom(rData); err != nil {
		return ret, err
	}

	*t = *NewRuneTrie(data, dict)

	return ret, nil
}

func (t *RuneTrie) GetDict() runes.Dict {
	return runes.Dict(t.dict)
}

func NewRuneTrieBuilder(opt ...da.Option) *RuneTrieBuilder {
	return &RuneTrieBuilder{
		builder: da.NewBuilder(opt...),
	}
}

func (b *RuneTrieBuilder) Build(ks doublearray.Walker, dict runes.Dict) (*RuneTrie, error) {
	data := da.New()
	if err := b.builder.Build(data, ks); err != nil {
		return nil, err
	}
	return NewRuneTrie(data, dict), nil
}

func (b *RuneTrieBuilder) BuildSlice(xs []string) (*RuneTrie, error) {
	dict := runes.New(xs)
	ks, err := b.keyset(xs, dict)
	if err != nil {
		return nil, err
	}
	return b.Build(ks, dict)
}

func (b *RuneTrieBuilder) BuildFromLines(r io.Reader) (*RuneTrie, error) {
	ks, dict, err := runes.FromWalker(token.NewLinedString(r))
	if err != nil {
		return nil, err
	}
	return b.Build(ks, dict)
}

func (*RuneTrieBuilder) keyset(ss []string, d runes.Dict) (keyset.KeySet, error) {
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
