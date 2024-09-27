package hairetsu

import (
	"io"
	"iter"

	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/codec/u16s"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
)

type doubleByteDict interface {
	codec.WordDict[uint16]
}
type doubleByteEncoder interface {
	codec.Encoder[[]byte]
}

var (
	_ doubleByteDict    = (u16s.WordDict)(nil)
	_ doubleByteEncoder = (*u16s.Encoder[u16s.Dict])(nil)
)

type DoubleByteTrie[T doubleByteDict] struct {
	data doublearray.Nodes
	dict T
}

type DoubleByteTrieBuilder[T doubleByteDict] struct {
	builder *doublearray.Builder
	dict    T
}

func NewDoubleByteTrie[T doubleByteDict](
	data doublearray.Nodes,
	dict T,
) *DoubleByteTrie[T] {
	return &DoubleByteTrie[T]{
		data: data,
		dict: dict,
	}
}

func (t *DoubleByteTrie[T]) InlineSearcher() *u16s.InlineSearcher[T, doublearray.Nodes] {
	return u16s.NewInlineSearcher(t.data, t.dict)
}
func (t *DoubleByteTrie[T]) Searcher() *codec.Searcher[[]byte, doublearray.Nodes] {
	return codec.NewSearcher(u16s.NewEncoder(t.dict), t.data)
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
	data := doublearray.New()
	ret, err := data.ReadFrom(r)
	n += ret
	if err == nil || err == io.EOF {
		t.data = data
	}
	return n, err
}

func NewDoubleByteTrieBuilder[T doubleByteDict](dict T, opt ...doublearray.Option) *DoubleByteTrieBuilder[T] {
	return &DoubleByteTrieBuilder[T]{
		builder: doublearray.NewBuilder(opt...),
		dict:    dict,
	}
}

func (b *DoubleByteTrieBuilder[T]) BuildFromLines(r io.ReadSeeker) (*DoubleByteTrie[T], error) {
	f := b.builder.Factory()
	err := u16s.FromReadSeeker(r, f, b.dict)
	if err != nil {
		return nil, err
	}
	return buildDoubleByteTrie(f, b.dict)
}

func buildDoubleByteTrie[T doubleByteDict](f *doublearray.Factory, dict T) (*DoubleByteTrie[T], error) {
	data, err := f.Done()
	if err != nil {
		return nil, err
	}
	return NewDoubleByteTrie(data, dict), nil
}
