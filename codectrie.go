package hairetsu

import (
	"io"
	"iter"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/codec/doublebyte"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type DoubleByteTrie struct {
	data da.Nodes
	dict doublebyte.Dict
}

type DoubleByteTrieBuilder struct {
	builder *da.Builder
}

func NewDoubleByteTrie(data da.Nodes, dict doublebyte.Dict) *DoubleByteTrie {
	return &DoubleByteTrie{
		data: data,
		dict: dict,
	}
}

func (t *DoubleByteTrie) Searcher() *codec.Searcher[[]byte] {
	return codec.NewSearcher(doublebyte.NewEncoder(t.dict), t.data)
}

func (t *DoubleByteTrie) Leafs() iter.Seq[item.Item] {
	return doublearray.Leafs(t.data)
}

func (t *DoubleByteTrie) WriteTo(w io.Writer) (int64, error) {
	ret, err := t.dict.WriteTo(w)
	if err != nil {
		return ret, err
	}
	n, err := t.data.WriteTo(w)
	ret += n
	return ret, err
}

func (t *DoubleByteTrie) ReadFrom(r io.Reader) (int64, error) {
	dict := doublebyte.ArrayDict{}
	n, err := dict.ReadFrom(r)
	if err != nil {
		return n, err
	}
	data := da.New()
	ret, err := data.ReadFrom(r)
	n += ret
	if err == nil || err == io.EOF {
		t.dict = dict
		t.data = data
	}
	return n, err
}

func NewDoubleByteTrieBuilder(opt ...da.Option) *DoubleByteTrieBuilder {
	return &DoubleByteTrieBuilder{
		builder: da.NewBuilder(opt...),
	}
}

func (b *DoubleByteTrieBuilder) BuildFromLines(r io.ReadSeeker) (*DoubleByteTrie, error) {
	f := b.builder.Factory()
	dict, err := doublebyte.FromReadSeeker(r, f)
	if err != nil {
		return nil, err
	}
	return buildDoubleByteTrie(f, dict)
}

func buildDoubleByteTrie(f *da.Factory, dict doublebyte.Dict) (*DoubleByteTrie, error) {
	data, err := f.Done()
	if err != nil {
		return nil, err
	}
	return NewDoubleByteTrie(data, dict), nil
}
