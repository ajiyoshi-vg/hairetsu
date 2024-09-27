package codec

import (
	"io"

	"github.com/ajiyoshi-vg/hairetsu/word"
	"golang.org/x/exp/constraints"
)

var (
	_ WordDict[int]        = (*Identity[int])(nil)
	_ Dict[word.Code, int] = (*InverseIdentity[int])(nil)
)

type Identity[T constraints.Integer] struct{}

func (*Identity[T]) Code(x T) word.Code {
	return word.Code(x)
}
func (*Identity[T]) Inverse() Dict[word.Code, T] {
	return &InverseIdentity[T]{}
}
func (*Identity[T]) Fill(count map[T]int) {
}
func (*Identity[T]) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}
func (*Identity[T]) ReadFrom(r io.Reader) (int64, error) {
	return 0, nil
}

type InverseIdentity[T constraints.Integer] struct{}

func (*InverseIdentity[T]) Code(x word.Code) T {
	return T(x)
}
func (*InverseIdentity[T]) Inverse() Dict[T, word.Code] {
	return &Identity[T]{}
}
