package hairetsu

import (
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/stretchr/testify/assert"
)

func BenchmarkMmap(b *testing.B) {
	trie, err := readIndex("byte.dat")
	assert.NoError(b, err)

	mmap, err := doublearray.NewMmap("byte.dat")
	assert.NoError(b, err)

	b.Run("trie", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range ws {
				if id, err := trie.ExactMatchSearch(v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	b.Run("mmap", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range ws {
				if id, err := mmap.ExactMatchSearch(v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	b.Run("mmap+i", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range ws {
				if id, err := doublearray.ExactMatchSearchInterface(mmap, v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
	b.Run("mmap+p", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range ws {
				if id, err := doublearray.ExactMatchSearchPointerMmap(mmap, v); err != nil {
					b.Fatalf("unexpected error, missing a keyword %v, id=%v, err=%v", v, id, err)
				}
			}
		}
	})
}
