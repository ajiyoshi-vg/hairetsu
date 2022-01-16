# hairetsu

hairetsu is a TRIE implementation by double array.

**alpha quality** : things would change.

## feature

* support ExactMatchSearch CommonPrefixSearch
* can use any binary as a label
  * including `\0`
  * Some other implementations treat `\0` as the end of string, so that they can't use keys including `\0` as a label.
* can customize the character code of the label
* can use as a key value store
  * can store 30bit uint value as a leaf

## how to use

build byte based TRIE and query

```go
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
```

build rune based trie and query

```go
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
```

## test

```bash
$ make test
```
