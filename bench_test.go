package hairetsu

import (
	"bufio"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/runedict"
	"github.com/stretchr/testify/assert"
)

func BenchmarkByteSearch(b *testing.B) {
	da, err := readIndex("byte.dat")
	assert.NoError(b, err)
	log.Println(da.Stat())

	trie := NewByteTrie(da)

	bs, err := readByteLines("bench.dat")
	assert.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		idx := node.Index(rand.Intn(len(bs)))
		actual, err := trie.ExactMatchSearch(bs[idx])
		assert.NoError(b, err)
		assert.Equal(b, int(idx), int(actual))
	}
}

func BenchmarkRuneSearch(b *testing.B) {
	da, err := readIndex("rune.dat")
	assert.NoError(b, err)
	log.Println(da.Stat())

	dict, err := readDict("rune.dat.dict")
	assert.NoError(b, err)

	trie := NewRuneTrie(da, dict)

	bs, err := readRuneLines("bench.dat")
	assert.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		idx := node.Index(rand.Intn(len(bs)))
		actual, err := trie.ExactMatchSearch(bs[idx])
		assert.NoError(b, err)
		assert.Equal(b, int(idx), int(actual))
	}
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
