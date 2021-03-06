package hairetsu

import (
	"bufio"
	"os"
	"testing"
	"time"

	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/overhead"
	"github.com/ajiyoshi-vg/hairetsu/token"
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
	file, err := os.Open("head.dat")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = token.NewLinedString(file).Walk(func(x string) error {
		bs = append(bs, []byte(x))
		ws = append(ws, word.FromBytes([]byte(x)))
		ss = append(ss, x)
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func BenchmarkTrie(b *testing.B) {
	b.Run("dartsclone", func(b *testing.B) {
		trie, err := dartsclone.Open("darts.trie")
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
		trie, err := readIndex("byte.trie")
		assert.NoError(b, err)
		b.Logf("byte.trie:%s", doublearray.GetStat(trie))
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
		da, err := readIndex("byte.trie")
		b.Logf("byte.trie:%s", doublearray.GetStat(da))
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
		trie := NewRuneTrie(nil, nil)
		{
			file, err := os.Open("rune.trie")
			assert.NoError(b, err)
			defer file.Close()

			_, err = trie.ReadFrom(bufio.NewReader(file))
			assert.NoError(b, err)
			b.Logf("rune.trie:%s", doublearray.GetStat(trie.data))
		}
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
	b.Run("dict", func(b *testing.B) {
		var trie DictTrie
		{
			file, err := os.Open("dict.trie")
			assert.NoError(b, err)
			defer file.Close()

			_, err = trie.ReadFrom(bufio.NewReader(file))
			if err != nil {
				b.Fatal(err)
				assert.NoError(b, err)
			}
			b.Logf("dict.trie:%s", doublearray.GetStat(trie.data))
		}
		b.Run("exact", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range bs {
					if id, err := trie.ExactMatchSearch(v); err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", string(v), id, err)
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
}
func BenchmarkOverhead(b *testing.B) {
	trie, err := readIndex("byte.trie")
	assert.NoError(b, err)

	start := time.Now()
	mmap, err := doublearray.OpenMmap("byte.trie")
	assert.NoError(b, err)
	b.Logf("OpenMmap %s", time.Since(start))

	b.Run("method", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range ws {
				if id, err := trie.ExactMatchSearch(v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	b.Run("pointer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range bs {
				if id, err := overhead.ExactMatchSearchPointer(trie, v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	b.Run("interface", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range bs {
				if id, err := overhead.ExactMatchSearchInterface(trie, v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	b.Run("mmap-m", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range ws {
				if id, err := mmap.ExactMatchSearch(v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	b.Run("mmap-p", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range bs {
				if id, err := overhead.ExactMatchSearchPointerMmap(mmap, v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	b.Run("mmap-i", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range bs {
				if id, err := overhead.ExactMatchSearchInterface(mmap, v); err != nil {
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
	_, err = da.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return da, nil
}
