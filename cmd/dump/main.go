package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ajiyoshi-vg/hairetsu"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/schollz/progressbar"
)

type option struct {
	in   string
	out  string
	kind string
}

var opt option

func init() {
	flag.StringVar(&opt.in, "in", "bench.dat", "line sep text default: bench.txt")
	flag.StringVar(&opt.out, "o", "out.dat", "output")
	flag.StringVar(&opt.kind, "kind", "byte", "[rune|byte] default: byte")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	switch opt.kind {
	case "byte":
		return dumpByte()
	case "rune":
		return dumpRune()
	default:
		return fmt.Errorf("unkown kind %s", opt.kind)
	}
}

func dumpByte() error {
	p := doublearray.OptionProgress(progressbar.New(0))
	trie, err := hairetsu.NewByteTrieBuilder(p).BuildFromFile(opt.in)
	if err != nil {
		return err
	}
	return dumpDoubleArray(trie, opt.out)
}

func dumpRune() error {
	p := doublearray.OptionProgress(progressbar.New(0))
	trie, err := hairetsu.NewRuneTrieBuilder(p).BuildFromFile(opt.in)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s.dict", opt.out)
	if err := dumpRuneDict(trie.GetDict(), path); err != nil {
		return err
	}
	return dumpDoubleArray(trie, opt.out)
}

type writable interface {
	WriteTo(io.Writer) (int64, error)
}

func dumpDoubleArray(data writable, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	w := bufio.NewWriter(out)
	defer w.Flush()
	if _, err := data.WriteTo(w); err != nil {
		return err
	}
	return nil
}

func dumpRuneDict(dict hairetsu.Dict, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	w := bufio.NewWriter(out)
	defer w.Flush()

	buf, err := dict.MarshalText()
	if err != nil {
		return err
	}
	_, err = w.WriteString(buf)
	return err
}
