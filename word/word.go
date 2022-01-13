package word

import (
	"fmt"
	"math"
	"sort"
)

const (
	EOS  = Code(0)
	NONE = Code(math.MaxUint32)
)

type Code uint32
type Word []Code

func (x Word) At(i int) Code {
	if i < len(x) {
		return x[i]
	}
	return EOS
}
func (x Word) Bytes() ([]byte, error) {
	ret := make([]byte, 0, len(x))
	for _, b := range x {
		if b > math.MaxUint8+1 {
			return nil, fmt.Errorf("bad code(%d) > MaxUint8", b)
		}
		ret = append(ret, byte(b-1))
	}
	return ret, nil
}
func (x Word) Sort() {
	sort.Slice(x, func(i, j int) bool {
		return x[i]-x[j] < 0
	})
}

func FromBytes(data []byte) Word {
	ret := make(Word, 0, len(data))
	for _, b := range data {
		ret = append(ret, Code(b)+1)
	}
	return ret
}

func Sort(data []Word) {
	sort.Slice(data, func(i, j int) bool {
		return Compare(data[i], data[j]) < 0
	})
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
