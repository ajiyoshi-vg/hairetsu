package example

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/stretchr/testify/assert"
)

func TestByteTrie(t *testing.T) {
	data := [][]byte{
		[]byte("aa"),
		[]byte("aaa"),
		[]byte("ab"),
		[]byte("abb"),
		[]byte("abc"),
		[]byte("abcd"),
		[]byte("b"),
		[]byte("ba"),
		[]byte("bb"),
		[]byte("c"),
		[]byte("cd"),
		[]byte("cddd"),
		[]byte("ccd"),
		[]byte("ddd"),
		[]byte("eab"),
		[]byte("日本語"),
		[]byte{math.MaxUint8, 0, math.MaxInt8},
	}

	trie, err := hairetsu.NewByteTrieBuilder().BuildSlice(data)
	assert.NoError(t, err)

	for i, x := range data {
		actual, err := trie.ExactMatchSearch(x)
		assert.NoError(t, err, x)
		assert.Equal(t, node.Index(i), actual)
	}

	ng := [][]byte{
		[]byte("a"),
		[]byte("aac"),
	}

	for _, x := range ng {
		_, err := trie.ExactMatchSearch(x)
		assert.Error(t, err, x)
	}

	target := []byte("abcedfg")
	is, err := trie.CommonPrefixSearch(target)
	assert.NoError(t, err)

	n := 0
	for _, x := range data {
		if bytes.HasPrefix(target, x) {
			n++
		}
	}

	assert.Equal(t, n, len(is))
}

func TestRuneTrie(t *testing.T) {
	data := []string{
		"aa",
		"aaa",
		"ab",
		"abb",
		"abc",
		"abcd",
		"b",
		"ba",
		"bb",
		"c",
		"cd",
		"cddd",
		"ccd",
		"ddd",
		"eab",
		"日本語",
	}

	trie, err := hairetsu.NewRuneTrieBuilder().BuildSlice(data)
	assert.NoError(t, err)

	for i, x := range data {
		actual, err := trie.ExactMatchSearch(x)
		assert.NoError(t, err, x)
		assert.Equal(t, node.Index(i), actual)
	}

	ng := []string{
		"a",
		"aac",
	}

	for _, x := range ng {
		_, err := trie.ExactMatchSearch(x)
		assert.Error(t, err, x)
	}

	target := "abcedfg"
	is, err := trie.CommonPrefixSearch(target)
	assert.NoError(t, err)

	n := 0
	for _, x := range data {
		if strings.HasPrefix(target, x) {
			n++
		}
	}

	assert.Equal(t, n, len(is))
}
