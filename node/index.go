package node

import (
	"math/rand"
	"reflect"

	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Index uint32

const MaxIndex = 1<<29 - 1

func (x Index) Forward(c word.Code) Index {
	switch c {
	case word.EOS:
		return x ^ 0
	case word.Separator:
		return x ^ 1
	default:
		return x ^ Index(c+word.ReservedCodes)
	}
}

// Backword : returns x such that x == offset.Forward(c)
func (x Index) Backward(c word.Code) Index {
	switch c {
	case word.EOS:
		return x ^ 0
	case word.Separator:
		return x ^ 1
	default:
		return x ^ Index(c+word.ReservedCodes)
	}
}

// Label : returns c such that offset.Forward(c) == child
func Label(offset, child Index) word.Code {
	c := word.Code(offset ^ child)
	switch c {
	case 0:
		return word.EOS
	case 1:
		return word.Separator
	default:
		return c - word.ReservedCodes
	}
}

func (x Index) Generate(r *rand.Rand, size int) reflect.Value {
	ret := Index(rand.Int31n(int32(MaxIndex)))
	return reflect.ValueOf(ret)
}
