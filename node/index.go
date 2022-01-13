package node

import (
	"math/rand"
	"reflect"

	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Index uint32

const MaxIndex = 1<<30 - 1

func (x Index) Forward(c word.Code) Index {
	return x ^ Index(c)
}

// Backword - return x such that x == offset.Forward(c)
func (x Index) Backward(c word.Code) (Index, error) {
	offset := x ^ Index(c)
	return offset, nil
}

func (x Index) Generate(r *rand.Rand, size int) reflect.Value {
	ret := Index(rand.Int31n(int32(MaxIndex)))
	return reflect.ValueOf(ret)
}
