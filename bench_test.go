package hairetsu

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/runedict"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/ikawaha/dartsclone"
	"github.com/stretchr/testify/assert"
)

var (
	ss []string
	bs [][]byte
	ws []word.Word
)

func init() {
	var err error
	bs, err = readByteLines("head.dat")
	if err != nil {
		panic(err)
	}

	ss = make([]string, 0, len(bs))
	ws = make([]word.Word, 0, len(bs))
	for _, b := range bs {
		ws = append(ws, word.FromBytes(b))
		ss = append(ss, string(b))
	}
}

func BenchmarkTrie(b *testing.B) {
	b.Run("dartsclone", func(b *testing.B) {
		trie, err := dartsclone.Open("darts.dat")
		assert.NoError(b, err)
		b.Run("exact", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range ss {
					if id, _, err := trie.ExactMatchSearch(v); id < 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
					}
				}
			}
		})
		b.Run("prefix", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range ss {
					if ret, err := trie.CommonPrefixSearch(v, 0); len(ret) == 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, err=%v", v, err)
					}
				}
			}
		})
	})
	b.Run("da", func(b *testing.B) {
		trie, err := readIndex("byte.dat")
		assert.NoError(b, err)
		b.Run("exact", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range ws {
					if id, err := trie.ExactMatchSearch(v); err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
					}
				}
			}
		})
		b.Run("prefix", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range ws {
					if ret, err := trie.CommonPrefixSearch(v); len(ret) == 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, err=%v", v, err)
					}
				}
			}
		})
	})
	b.Run("byte", func(b *testing.B) {
		da, err := readIndex("byte.dat")
		assert.NoError(b, err)
		trie := NewByteTrie(da)
		b.Run("exact", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range bs {
					if id, err := trie.ExactMatchSearch(v); err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
					}
				}
			}
		})
		b.Run("prefix", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range bs {
					if ret, err := trie.CommonPrefixSearch(v); len(ret) == 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, err=%v", v, err)
					}
				}
			}
		})
	})
	b.Run("rune", func(b *testing.B) {
		da, err := readIndex("rune.dat")
		assert.NoError(b, err)
		dict, err := readDict("rune.dat.dict")
		assert.NoError(b, err)
		trie := NewRuneTrie(da, dict)
		b.Run("exact", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range ss {
					if id, err := trie.ExactMatchSearch(v); err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
					}
				}
			}
		})
		b.Run("prefix", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range ss {
					if ret, err := trie.CommonPrefixSearch(v); len(ret) == 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, err=%v", v, err)
					}
				}
			}
		})
	})
}
func BenchmarkOverhead(b *testing.B) {
	trie, err := readIndex("byte.dat")
	assert.NoError(b, err)
	b.Run("method", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range ws {
				if id, err := trie.ExactMatchSearch(v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	/* need go 1.18
	b.Run("generics", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range ws {
				if id, err := exactMatchSearchGenerics(trie, v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	*/
	b.Run("pointer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range ws {
				if id, err := doublearray.ExactMatchSearchPointer(trie, v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	b.Run("interface", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range ws {
				if id, err := doublearray.ExactMatchSearchInterface(trie, v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
}

func readIndex(path string) (*doublearray.DoubleArray, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := bufio.NewReader(file)
	da := doublearray.New()
	_, err = doublearray.NewBuilder().ReadFrom(da, r)
	if err != nil {
		return nil, err
	}
	return da, nil
}
func readDict(path string) (runedict.RuneDict, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	ret := runedict.RuneDict{}
	if err := ret.UnmarshalText(string(buf)); err != nil {
		return nil, err
	}
	return ret, nil
}

func readRuneLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := make([]string, 0, 1000)
	scan := bufio.NewScanner(file)
	for i := 0; scan.Scan(); i++ {
		line := scan.Text()
		ret = append(ret, line)
	}

	return ret, nil
}
func readByteLines(path string) ([][]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := make([][]byte, 0, 1000)
	scan := bufio.NewScanner(file)
	for i := 0; scan.Scan(); i++ {
		line := scan.Text()
		ret = append(ret, []byte(line))
	}

	return ret, nil
}
