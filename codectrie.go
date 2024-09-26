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

type DoubleByteTrie[Dict doublebyte.Dict] struct {
	data da.Nodes
	dict Dict
}

type DoubleByteTrieBuilder[Dict doublebyte.Dict] struct {
	builder *da.Builder
	dict    Dict
}

func NewDoubleByteTrie[Dict doublebyte.Dict](
	data da.Nodes,
	dict Dict,
) *DoubleByteTrie[Dict] {
	return &DoubleByteTrie[Dict]{
		data: data,
		dict: dict,
	}
}

func (t *DoubleByteTrie[Dict]) InlineSearcher() *doublebyte.InlineSearcher[Dict] {
	return doublebyte.NewInlineSearcher(t.data, t.dict)
}
func (t *DoubleByteTrie[Dict]) Searcher() *codec.Searcher[[]byte] {
	return codec.NewSearcher(doublebyte.NewEncoder(t.dict), t.data)
}

func (t *DoubleByteTrie[Dict]) Leafs() iter.Seq[item.Item] {
	return doublearray.Leafs(t.data)
}

func (t *DoubleByteTrie[Dict]) WriteTo(w io.Writer) (int64, error) {
	ret, err := t.dict.WriteTo(w)
	if err != nil {
		return ret, err
	}
	n, err := t.data.WriteTo(w)
	ret += n
	return ret, err
}

func (t *DoubleByteTrie[Dict]) ReadFrom(r io.Reader) (int64, error) {
	n, err := t.dict.ReadFrom(r)
	if err != nil {
		return n, err
	}
	data := da.New()
	ret, err := data.ReadFrom(r)
	n += ret
	if err == nil || err == io.EOF {
		t.data = data
	}
	return n, err
}

func NewDoubleByteTrieBuilder[D doublebyte.Dict](dict D, opt ...da.Option) *DoubleByteTrieBuilder[D] {
	return &DoubleByteTrieBuilder[D]{
		builder: da.NewBuilder(opt...),
		dict:    dict,
	}
}

func (b *DoubleByteTrieBuilder[D]) BuildFromLines(r io.ReadSeeker) (*DoubleByteTrie[D], error) {
	f := b.builder.Factory()
	dict, err := doublebyte.FromReadSeeker(b.dict, r, f)
	if err != nil {
		return nil, err
	}
	return buildDoubleByteTrie(f, dict)
}

func buildDoubleByteTrie[D doublebyte.Dict](f *da.Factory, dict D) (*DoubleByteTrie[D], error) {
	data, err := f.Done()
	if err != nil {
		return nil, err
	}
	return NewDoubleByteTrie(data, dict), nil
}
