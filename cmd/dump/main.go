package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/keytree"
	"github.com/ajiyoshi-vg/hairetsu/node"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/pkg/profile"
	"github.com/schollz/progressbar"
)

type option struct {
	in       string
	out      string
	validate bool
}

var opt option

func init() {
	flag.StringVar(&opt.in, "in", "bench.dat", "line sep text default: bench.txt")
	flag.StringVar(&opt.out, "o", "trie.dat", "output")
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
	ks, err := readFile(opt.in)
	if err != nil {
		return err
	}
	da, err := build(ks)
	if err != nil {
		return err
	}

	out, err := os.Create(opt.out)
	if err != nil {
		return err
	}
	defer out.Close()
	w := bufio.NewWriter(out)
	defer w.Flush()
	if _, err := da.WriteTo(w); err != nil {
		return err
	}
	return nil
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

func build(data doublearray.Walker) (*doublearray.DoubleArray, error) {
	x := da.New()
	if err := da.NewBuilder(da.OptionProgress(progressbar.New(1))).Build(x, data); err != nil {
		return nil, err
	}
	log.Println("build finished")
	log.Println(x.Stat())
	err := data.WalkLeaf(func(key word.Word, val uint32) error {
		actual, err := x.ExactMatchSearch(key)
		if err != nil {
			return err
		}
		if actual != node.Index(val) {
			return fmt.Errorf("search key(%v): want %d got %d", key, val, actual)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return x, nil
}
