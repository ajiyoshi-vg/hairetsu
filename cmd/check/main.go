package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/keytree"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/pkg/profile"
)

type option struct {
	in   string
	data string
}

var opt option

func init() {
	flag.StringVar(&opt.in, "in", "bench.dat", "line sep text default: bench.txt")
	flag.StringVar(&opt.data, "data", "trie.dat", "dumped file")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	defer profile.Start(profile.ProfilePath(".")).Stop()
	da, err := load()
	if err != nil {
		return err
	}

	ks, err := readFile(opt.in)
	if err != nil {
		return err
	}

	return ks.WalkLeaf(func(key word.Word, val uint32) error {
		actual, err := da.ExactMatchSearch(key)
		if err != nil {
			return err
		}
		expect := node.Index(val)
		if actual != expect {
			log.Printf("search(%v) want %d got %d", key, expect, actual)
		}
		return nil
	})
}

func readFile(path string) (*keytree.Tree, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := keytree.New()
	scan := bufio.NewScanner(file)
	for i := 0; scan.Scan(); i++ {
		line := scan.Text()
		ret.Put(word.FromBytes([]byte(line)), uint32(i))
	}
	return ret, nil
}

func load() (*doublearray.DoubleArray, error) {
	file, err := os.Open(opt.data)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r := bufio.NewReader(file)

	x := da.New()
	if _, err := da.NewBuilder().ReadFrom(x, r); err != nil {
		return nil, err
	}
	log.Println("load finished")
	log.Println(x.Stat())
	return x, nil
}
