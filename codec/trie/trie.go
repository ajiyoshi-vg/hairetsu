package trie

import (
	"bufio"
	"io"
)

func multiCopy(w io.Writer, xs ...io.WriterTo) (int64, error) {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	var ret int64
	for _, x := range xs {
		n, err := x.WriteTo(bw)
		ret += n
		if err != nil {
			return ret, err
		}
	}

	return ret, nil
}
