package word

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
)

/*
[math.MaxUint32 - ReservedCodes ... math.MaxUint32] are reserved special code.

Unknown is invalid code. Words must not contain Unknown code.
Separator is The separetor. Users can use it for separated namespace in key.
EOS is the end of string. Words must not contain EOS code.
*/
const ReservedCodes = 8
const (
	Unknown Code = iota + math.MaxUint32 - ReservedCodes
	reserved1
	reserved2
	reserved3
	reserved4
	reserved5
	Separator
	EOS
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

func FromString(s string) Word {
	ret := make(Word, 0, len(s))
	for _, b := range s {
		ret = append(ret, Code(b))
	}
	return ret
}

func Compare(a, b Word) int {
	for i, x := range a {
		if i >= len(b) {
			return 1
		}
		if x < b[i] {
			return -1
		}
		if x > b[i] {
			return 1
		}
	}
	if len(a) < len(b) {
		return -1
	}
	return 0
}

func Equal(a, b Word) bool {
	return Compare(a, b) == 0
}

func WithNameSpace(ns, key []byte) Word {
	ret := make(Word, 0, len(ns)+len(key)+1)
	for _, b := range ns {
		ret = append(ret, Code(b))
	}
	ret = append(ret, Separator)
	for _, b := range key {
		ret = append(ret, Code(b))
	}
	return ret
}

// UnNameSpace : inverse of WithNameSpace(). it's for test purpose
func (x Word) UnNameSpace() ([]byte, []byte, error) {
	var head, tail []byte
	separated := false
	for _, c := range x {
		if c == Separator {
			separated = true
			continue
		}
		if c > math.MaxUint8 {
			return nil, nil, fmt.Errorf("bad %v code(%d) %d", x, c, Separator)
		}
		if separated {
			tail = append(tail, byte(c))
		} else {
			head = append(head, byte(c))
		}
	}
	if separated {
		return head, tail, nil
	}
	return nil, head, nil
}

func (x Code) Generate(r *rand.Rand, size int) reflect.Value {
	ret := Code(rand.Uint32())
	return reflect.ValueOf(ret)
}

func Reverse(w Word) {
	i, j := 0, len(w)-1
	for i < j {
		w[i], w[j] = w[j], w[i]
		i++
		j--
	}
}
