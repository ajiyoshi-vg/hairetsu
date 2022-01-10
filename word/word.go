package word

import "math"

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

func FromByte(data []byte) Word {
	ret := make(Word, 0, len(data))
	for _, b := range data {
		ret = append(ret, Code(b+1))
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
