package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/ajiyoshi-vg/hairetsu"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/keytree"
	"github.com/ajiyoshi-vg/hairetsu/lines"
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
	start := time.Now()
	log.Printf("%s -> %s start", opt.in, opt.out)
	if err := run(); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s -> %s dumped in %s", opt.in, opt.out, time.Since(start))
	os.Exit(0)
}

func run() error {
	defer fmt.Println("finish")
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
	file, err := os.Open(opt.in)
	if err != nil {
		return err
	}
	defer file.Close()
	ks, err := keytree.FromLines(file)
	if err != nil {
		return err
	}
	p := doublearray.OptionProgress(progressbar.New(0))
	trie, err := hairetsu.NewByteTrieBuilder(p).Build(ks)
	if err != nil {
		return err
	}
	return writeTo(trie, opt.out)
}

func dumpRune() error {
	file, err := os.Open(opt.in)
	if err != nil {
		return err
	}
	defer file.Close()

	p := doublearray.OptionProgress(progressbar.New(0))
	trie, err := hairetsu.NewRuneTrieBuilder(p).BuildFromLines(file)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s.dict", opt.out)
	if err := writeTo(trie.GetDict(), path); err != nil {
		return err
	}
	return writeTo(trie, opt.out)
}

func dumpDarts() error {
	file, err := os.Open(opt.in)
	if err != nil {
		return err
	}
	defer file.Close()
	ss, err := lines.StringSlice(file)
	if err != nil {
		return err
	}

	p := dartsprog.New()
	p.SetMaximum(len(ss))
	b := dartsclone.NewBuilder(p)
	if err := b.Build(ss, nil); err != nil {
		return err
	}

	return writeTo(b, opt.out)
}

func writeTo(data io.WriterTo, path string) error {
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
