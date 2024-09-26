package hairetsu

import (
	"bytes"
	"encoding/binary"
	"io"

	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/runes"
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
	if err != nil {
		return 0, err
	}

	return io.Copy(w, io.MultiReader(size, dict, data))
}

func (t *RuneTrie) ReadFrom(r io.Reader) (int64, error) {
	buf, err := io.ReadAll(r)
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
	data := da.New()
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

func (b *RuneTrieBuilder) BuildFromSlice(xs []string) (*RuneTrie, error) {
	f := b.builder.Factory()
	dict, err := runes.FromSlice(xs, f)
	if err != nil {
		return nil, err
	}
	return buildRuneTrie(f, dict)
}

func (b *RuneTrieBuilder) BuildFromLines(r io.ReadSeeker) (*RuneTrie, error) {
	f := b.builder.Factory()
	dict, err := runes.FromReader(r, f)
	if err != nil {
		return nil, err
	}
	return buildRuneTrie(f, dict)
}

func buildRuneTrie(f *da.Factory, dict runes.Dict) (*RuneTrie, error) {
	trie, err := f.Done()
	if err != nil {
		return nil, err
	}
	return NewRuneTrie(trie, dict), nil
}
