package word

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
)

/*
EOS is the end of string. Words must not contain EOS code.
NONE is invalid code. Words must not contain NONE code.
SEP is The separetor. Users can use it for separated namespace in key.
*/
const (
	EOS  = Code(math.MaxUint32)
	SEP  = Code(math.MaxUint32 - 1)
	NONE = Code(math.MaxUint32 - 2)
)

type Code uint32
type Word []Code

func (x Word) At(i int) Code {
	if i < len(x) {
		return x[i]
	}
	return EOS
}

// Bytes : inverse of FromBytes(). it's for test purpose
func (x Word) Bytes() ([]byte, error) {
	ret := make([]byte, 0, len(x))
	for _, b := range x {
		if b > math.MaxUint8+1 {
			return nil, fmt.Errorf("bad code(%d) > MaxUint8", b)
		}
		ret = append(ret, byte(b))
	}
	return ret, nil
}

func FromBytes(data []byte) Word {
	ret := make(Word, 0, len(data))
	for _, b := range data {
		ret = append(ret, Code(b))
	}
	return ret
}

func WithNameSpace(ns, key []byte) Word {
	ret := make(Word, 0, len(ns)+len(key)+1)
	for _, b := range ns {
		ret = append(ret, Code(b))
	}
	ret = append(ret, SEP)
	for _, b := range key {
		ret = append(ret, Code(b))
	}
	return ret
}

func Compare(lhs, rhs Word) int {
	shorter := len(lhs)
	if len(rhs) < shorter {
		shorter = len(rhs)
	}

	for i := 0; i < shorter; i++ {
		if lhs[i] != rhs[i] {
			return int(lhs[i]) - int(rhs[i])
		}
	}

	return len(lhs) - len(rhs)
}

func (x Code) Generate(r *rand.Rand, size int) reflect.Value {
	ret := Code(rand.Uint32())
	return reflect.ValueOf(ret)
}
