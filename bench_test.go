package hairetsu

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu/codec"
	"github.com/ajiyoshi-vg/hairetsu/codec/bytes"
	"github.com/ajiyoshi-vg/hairetsu/codec/runes"
	"github.com/ajiyoshi-vg/hairetsu/codec/trie"
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
}

func BenchmarkCodec(b *testing.B) {
	byteEncoder := map[string]codec.Encoder[[]byte]{
		"u16s-a":  u16s.NewEncoder(u16s.NewArrayDict()),
		"u16s-i":  u16s.NewEncoder(u16s.NewIdentityDict()),
		"bytes-a": bytes.NewEncoder(bytes.NewArrayDict()),
		"bytes-i": bytes.NewEncoder(bytes.NewIdentityDict()),
	}
	runeEncoder := map[string]codec.Encoder[string]{
		"runes-m": runes.NewEncoder(runes.NewMapDict()),
		"runes-i": runes.NewEncoder(runes.NewIdentityDict()),
	}

	b.Run("file", func(b *testing.B) {
		keys := slices.Collect(maps.Keys(byteEncoder))
		slices.Sort(keys)
		for _, kind := range keys {
			b.Run(kind, func(b *testing.B) {
				t, err := trie.OpenFile(fmt.Sprintf("%s.trie", kind), byteEncoder[kind])
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

		keys = slices.Collect(maps.Keys(runeEncoder))
		slices.Sort(keys)
		for _, kind := range keys {
			b.Run(kind, func(b *testing.B) {
				file, encoder := fmt.Sprintf("%s.trie", kind), runeEncoder[kind]
				t, err := trie.OpenFile(file, encoder)
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
	b.Run("mmap", func(b *testing.B) {
		keys := slices.Collect(maps.Keys(byteEncoder))
		slices.Sort(keys)
		for _, kind := range keys {
			b.Run(kind, func(b *testing.B) {
				file, encoder := fmt.Sprintf("%s.trie", kind), byteEncoder[kind]
				t, err := trie.OpenMmap(file, encoder)
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
		keys = slices.Collect(maps.Keys(runeEncoder))
		slices.Sort(keys)
		for _, kind := range keys {
			b.Run(kind, func(b *testing.B) {
				file, encoder := fmt.Sprintf("%s.trie", kind), runeEncoder[kind]
				t, err := trie.OpenMmap(file, encoder)
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
		b.Run("byte", func(b *testing.B) {
			da, err := doublearray.OpenFile("bytes-i.trie")
			assert.NoError(b, err)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for _, v := range bs {
					if id, err := trie.BytesExactMatchSearch(da, v); err != nil {
						b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
					}
				}
			}
		})
		b.Run("rune", func(b *testing.B) {
			da, err := doublearray.OpenFile("runes-i.trie")
			assert.NoError(b, err)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for _, v := range ss {
					if id, err := trie.RunsExactMatchSearch(da, v); err != nil {
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
	b.Run("inline", func(b *testing.B) {
		b.Run("file", func(b *testing.B) {
			t, err := trie.OpenInlineFile("u16s-a.trie", u16s.NewArrayDict())
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
		b.Run("mmap", func(b *testing.B) {
			t, err := trie.OpenInlineFile("u16s-a.trie", u16s.NewArrayDict())
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
	})
	b.Run("word", func(b *testing.B) {
		t, err := doublearray.OpenFile("bytes-i.trie")
		assert.NoError(b, err)
		b.Logf("byte-i.trie:%s", doublearray.GetStat(t))
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
}
func BenchmarkOverhead(b *testing.B) {
	t, err := doublearray.OpenFile("byte.trie")
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
