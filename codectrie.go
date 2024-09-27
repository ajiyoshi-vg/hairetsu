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

type DoubleByteTrie[T doublebyte.FullDict] struct {
	data da.Nodes
	dict T
}

type DoubleByteTrieBuilder[T doublebyte.FullDict] struct {
	builder *da.Builder
	dict    T
}

func NewDoubleByteTrie[T doublebyte.FullDict](
	data da.Nodes,
	dict T,
) *DoubleByteTrie[T] {
	return &DoubleByteTrie[T]{
		data: data,
		dict: dict,
	}
}

func (t *DoubleByteTrie[T]) InlineSearcher() *doublebyte.InlineSearcher[T] {
	return doublebyte.NewInlineSearcher(t.data, t.dict)
}
func (t *DoubleByteTrie[T]) Searcher() *codec.Searcher[[]byte] {
	return codec.NewSearcher(doublebyte.NewEncoder(t.dict), t.data)
}

func (t *DoubleByteTrie[T]) Leafs() iter.Seq[item.Item] {
	return doublearray.Leafs(t.data)
}

func (t *DoubleByteTrie[T]) WriteTo(w io.Writer) (int64, error) {
	ret, err := t.dict.WriteTo(w)
	if err != nil {
		return ret, err
	}
	n, err := t.data.WriteTo(w)
	ret += n
	return ret, err
}

func (t *DoubleByteTrie[T]) ReadFrom(r io.Reader) (int64, error) {
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

func NewDoubleByteTrieBuilder[T doublebyte.FullDict](dict T, opt ...da.Option) *DoubleByteTrieBuilder[T] {
	return &DoubleByteTrieBuilder[T]{
		builder: da.NewBuilder(opt...),
		dict:    dict,
	}
}

func (b *DoubleByteTrieBuilder[T]) BuildFromLines(r io.ReadSeeker) (*DoubleByteTrie[T], error) {
	f := b.builder.Factory()
	err := doublebyte.FromReadSeeker(r, f, b.dict)
	if err != nil {
		return nil, err
	}
	return buildDoubleByteTrie(f, b.dict)
}

func buildDoubleByteTrie[T doublebyte.FullDict](f *da.Factory, dict T) (*DoubleByteTrie[T], error) {
	data, err := f.Done()
	if err != nil {
		return nil, err
	}
	return NewDoubleByteTrie(data, dict), nil
}
