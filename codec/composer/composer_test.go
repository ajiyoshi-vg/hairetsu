package composer

import (
	stdbytes "bytes"
	"io"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/codec/bytes"
	"github.com/ajiyoshi-vg/hairetsu/codec/runes"
	"github.com/ajiyoshi-vg/hairetsu/codec/u16s"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/stretchr/testify/assert"
)

type composer[T any] interface {
	Compose(r io.ReadSeeker) (*Trie[T, *doublearray.DoubleArray], error)
}

func TestBytesCompose(t *testing.T) {
	cases := map[string]struct {
		data    string
		ok      [][]byte
		unknown [][]byte
	}{
		"simple": {
			data: "a\nb\nc\n",
			ok: [][]byte{
				[]byte("a"),
				[]byte("b"),
				[]byte("c"),
			},
			unknown: [][]byte{
				[]byte("d"),
			},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			kinds := map[string]composer[[]byte]{
				"int16/m": NewInt16(u16s.NewMapDict()),
				"int16/a": NewInt16(u16s.NewArrayDict()),
				"int16/i": NewInt16(u16s.NewIdentityDict()),
				"bytes/m": NewBytes(bytes.NewMapDict()),
				"bytes/a": NewBytes(bytes.NewArrayDict()),
				"bytes/i": NewBytes(bytes.NewIdentityDict()),
			}
			for kind, x := range kinds {

				t.Run(kind, func(t *testing.T) {
					in := stdbytes.NewReader([]byte(c.data))

					trie, err := x.Compose(in)
					assert.NoError(t, err)

					s := trie.Searcher()

					for _, ok := range c.ok {
						_, err := s.ExactMatchSearch(ok)
						assert.NoError(t, err)
					}

					for _, ng := range c.unknown {
						_, err := s.ExactMatchSearch(ng)
						assert.Error(t, err)
					}
				})
			}
		})
	}
}

func TestStringCompose(t *testing.T) {
	cases := map[string]struct {
		data    string
		ok      []string
		unknown []string
	}{
		"simple": {
			data:    "a\nb\nc\n",
			ok:      []string{"a", "b", "c"},
			unknown: []string{"d"},
		},
	}

	for title, c := range cases {
		t.Run(title, func(t *testing.T) {
			kinds := map[string]composer[string]{
				"runes/m": NewRunes(runes.NewMapDict()),
				"runes/i": NewRunes(runes.NewIdentityDict()),
			}
			for kind, x := range kinds {
				t.Run(kind, func(t *testing.T) {
					in := stdbytes.NewReader([]byte(c.data))

					trie, err := x.Compose(in)
					assert.NoError(t, err)

					s := trie.Searcher()

					for _, ok := range c.ok {
						_, err := s.ExactMatchSearch(ok)
						assert.NoError(t, err)
					}

					for _, ng := range c.unknown {
						_, err := s.ExactMatchSearch(ng)
						assert.Error(t, err)
					}
				})
			}
		})
	}
}
