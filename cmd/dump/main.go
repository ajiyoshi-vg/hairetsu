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
	"github.com/ajiyoshi-vg/hairetsu/runedict"
	"github.com/ikawaha/dartsclone"
	dartsprog "github.com/ikawaha/dartsclone/progressbar"
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
	flag.StringVar(&opt.kind, "kind", "byte", "[rune|byte|darts] default: byte")
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
	case "darts":
		return dumpDarts()
	default:
		return fmt.Errorf("unkown kind %s", opt.kind)
	}
}

func dumpByte() error {
	ss, err := readLinesByte(opt.in)
	if err != nil {
		return err
	}
	p := doublearray.OptionProgress(progressbar.New(0))
	trie, err := hairetsu.NewByteTrieBuilder(p).BuildSlice(ss)
	if err != nil {
		return err
	}
	return dumpDoubleArray(trie, opt.out)
}

func dumpRune() error {
	ss, err := readLines(opt.in)
	if err != nil {
		return err
	}
	p := doublearray.OptionProgress(progressbar.New(0))
	trie, err := hairetsu.NewRuneTrieBuilder(p).BuildSlice(ss)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s.dict", opt.out)
	if err := dumpRuneDict(trie.GetDict(), path); err != nil {
		return err
	}
	return dumpDoubleArray(trie, opt.out)
}

func dumpDarts() error {
	ss, err := readLines(opt.in)
	if err != nil {
		return err
	}

	p := dartsprog.New()
	p.SetMaximum(len(ss))
	b := dartsclone.NewBuilder(p)
	if err := b.Build(ss, nil); err != nil {
		return err
	}

	out, err := os.Create(opt.out)
	if err != nil {
		return err
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	defer w.Flush()
	_, err = b.WriteTo(w)
	return err
}

func readLines(path string) ([]string, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	ret := make([]string, 0, 1000)
	scan := bufio.NewScanner(r)
	for i := 0; scan.Scan(); i++ {
		line := scan.Text()
		ret = append(ret, line)
	}
	return ret, nil
}
func readLinesByte(path string) ([][]byte, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	ret := make([][]byte, 0, 1000)
	scan := bufio.NewScanner(r)
	for i := 0; scan.Scan(); i++ {
		line := scan.Text()
		ret = append(ret, []byte(line))
	}
	return ret, nil
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

func dumpRuneDict(dict runedict.RuneDict, path string) error {
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
