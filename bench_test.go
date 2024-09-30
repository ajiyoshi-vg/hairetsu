package hairetsu

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/codec/bytes"
	"github.com/ajiyoshi-vg/hairetsu/codec/composer"
	"github.com/ajiyoshi-vg/hairetsu/codec/runes"
	"github.com/ajiyoshi-vg/hairetsu/codec/u16s"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/overhead"
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
	for x := range scan.ByteLines(file) {
		bs = append(bs, x)
		ws = append(ws, word.FromBytes(x))
		ss = append(ss, string(x))
	}
}

func BenchmarkTrie(b *testing.B) {
	b.Run("dartsclone", func(b *testing.B) {
		t, err := dartsclone.Open("darts.trie")
		assert.NoError(b, err)
		b.Run("exact", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range ss {
					if id, _, err := t.ExactMatchSearch(v); id < 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
					}
				}
			}
		})
		b.Run("prefix", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range ss {
					if ret, err := t.CommonPrefixSearch(v, 0); len(ret) == 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, err=%v", v, err)
					}
				}
			}
		})
	})
	b.Run("da", func(b *testing.B) {
		t, err := readIndex("byte.trie")
		assert.NoError(b, err)
		b.Logf("byte.trie:%s", doublearray.GetStat(t))
		b.Run("exact", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range ws {
					if id, err := t.ExactMatchSearch(v); err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
					}
				}
			}
		})
		b.Run("prefix", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range ws {
					if ret, err := t.CommonPrefixSearch(v); len(ret) == 0 || err != nil {
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
		t := NewByteTrie(da)
		b.Run("exact", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range bs {
					if id, err := t.ExactMatchSearch(v); err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
					}
				}
			}
		})
		b.Run("prefix", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range bs {
					if ret, err := t.CommonPrefixSearch(v); len(ret) == 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, err=%v", v, err)
					}
				}
			}
		})
	})
	b.Run("rune", func(b *testing.B) {
		t := NewRuneTrie(nil, nil)
		{
			file, err := os.Open("rune.trie")
			assert.NoError(b, err)
			defer file.Close()

			_, err = t.ReadFrom(bufio.NewReader(file))
			assert.NoError(b, err)
			b.Logf("rune.trie:%s", doublearray.GetStat(t.data))
		}
		b.Run("exact", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range ss {
					if id, err := t.ExactMatchSearch(v); err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
					}
				}
			}
		})
		b.Run("prefix", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range ss {
					if ret, err := t.CommonPrefixSearch(v); len(ret) == 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, err=%v", v, err)
					}
				}
			}
		})
	})
	b.Run("dict", func(b *testing.B) {
		var t DictTrie
		{
			file, err := os.Open("dict.trie")
			assert.NoError(b, err)
			defer file.Close()

			_, err = t.ReadFrom(bufio.NewReader(file))
			if err != nil {
				b.Fatal(err)
				assert.NoError(b, err)
			}
			b.Logf("dict.trie:%s", doublearray.GetStat(t.data))
		}
		b.Run("exact", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range bs {
					if id, err := t.ExactMatchSearch(v); err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", string(v), id, err)
					}
				}
			}
		})
		b.Run("prefix", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, v := range bs {
					if ret, err := t.CommonPrefixSearch(v); len(ret) == 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, err=%v", v, err)
					}
				}
			}
		})
	})
}

func BenchmarkCodec(b *testing.B) {
	b.Run("byte", func(b *testing.B) {
		kinds := map[string]codec.Encoder[[]byte]{
			"u16s-m":  u16s.NewEncoder((u16s.NewMapDict())),
			"u16s-a":  u16s.NewEncoder((u16s.NewArrayDict())),
			"u16s-i":  u16s.NewEncoder((u16s.NewIdentityDict())),
			"bytes-m": bytes.NewEncoder(bytes.NewMapDict()),
			"bytes-a": bytes.NewEncoder(bytes.NewArrayDict()),
			"bytes-i": bytes.NewEncoder(bytes.NewIdentityDict()),
		}
		keys := slices.Collect(maps.Keys(kinds))
		slices.Sort(keys)
		for _, kind := range keys {
			b.Run(kind, func(b *testing.B) {
				t := composer.NewFileTrie(kinds[kind])
				err := t.Open(fmt.Sprintf("%s.trie", kind))
				assert.NoError(b, err)
				s := t.Searcher()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					for _, v := range bs {
						if id, err := s.ExactMatchSearch(v); err != nil {
							b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", string(v), id, err)
						}
					}
				}
			})
		}

	})
	b.Run("rune", func(b *testing.B) {
		kinds := map[string]codec.Encoder[string]{
			"runes-m": runes.NewEncoder(runes.NewMapDict()),
			"runes-i": runes.NewEncoder(runes.NewIdentityDict()),
		}
		keys := slices.Collect(maps.Keys(kinds))
		slices.Sort(keys)
		for _, kind := range keys {
			b.Run(kind, func(b *testing.B) {
				t := composer.NewFileTrie(kinds[kind])
				err := t.Open(fmt.Sprintf("%s.trie", kind))
				assert.NoError(b, err)
				s := t.Searcher()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					for _, v := range ss {
						if id, err := s.ExactMatchSearch(v); err != nil {
							b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
						}
					}
				}
			})
		}
	})
	b.Run("misc", func(b *testing.B) {
		b.Run("old", func(b *testing.B) {
			da, err := readIndex("bytes-i.trie")
			assert.NoError(b, err)
			s := NewByteTrie(da)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for _, v := range bs {
					if id, err := s.ExactMatchSearch(v); err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
					}
				}
			}
		})
		b.Run("inline", func(b *testing.B) {
			t := composer.NewFileInline(u16s.NewArrayDict())
			err := t.Open("u16s-a.trie")
			assert.NoError(b, err)
			s := t.Searcher()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for _, v := range bs {
					if id, err := s.ExactMatchSearch(v); err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
					}
				}
			}
		})
		b.Run("darts", func(b *testing.B) {
			t, err := dartsclone.Open("darts.trie")
			assert.NoError(b, err)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for _, v := range ss {
					if id, _, err := t.ExactMatchSearch(v); id < 0 || err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
					}
				}
			}
		})
	})
}
func BenchmarkOverhead(b *testing.B) {
	t, err := readIndex("byte.trie")
	assert.NoError(b, err)

	start := time.Now()
	mmap, err := doublearray.OpenMmap("byte.trie")
	assert.NoError(b, err)
	b.Logf("OpenMmap %s", time.Since(start))

	b.Run("method", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range ws {
				if id, err := t.ExactMatchSearch(v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	b.Run("pointer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range bs {
				if id, err := overhead.ExactMatchSearchPointer(t, v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	b.Run("interface", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range bs {
				if id, err := overhead.ExactMatchSearchInterface(t, v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	b.Run("generics", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range bs {
				if id, err := overhead.ExactMatchSearchGenerics(t, v); err != nil {
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
	b.Run("mmap-g", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range bs {
				if id, err := overhead.ExactMatchSearchGenerics(mmap, v); err != nil {
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
