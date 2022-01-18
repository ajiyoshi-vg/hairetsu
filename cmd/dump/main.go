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
	"github.com/ajiyoshi-vg/hairetsu/token"
	"github.com/ajiyoshi-vg/hairetsu/word"
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
	flag.StringVar(&opt.kind, "kind", "byte", "[rune|byte|dict|darts] default: byte")
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

	file, err := os.Open(opt.in)
	if err != nil {
		return err
	}
	defer file.Close()

	switch opt.kind {
	case "byte":
		return dumpByte(file)
	case "rune":
		return dumpRune(file)
	case "dict":
		return dumpDict(file)
	case "darts":
		return dumpDarts(file)
	default:
		return fmt.Errorf("unkown kind %s", opt.kind)
	}
}

func dumpByte(file io.Reader) error {
	ks, err := fromLines(file)
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

func dumpRune(file io.Reader) error {
	p := doublearray.OptionProgress(progressbar.New(0))
	trie, err := hairetsu.NewRuneTrieBuilder(p).BuildFromLines(file)
	if err != nil {
		return err
	}
	return writeTo(trie, opt.out)
}

func dumpDict(file io.Reader) error {
	p := doublearray.OptionProgress(progressbar.New(0))
	trie, err := hairetsu.NewDictTrieBuilder(p).BuildFromLines(file)
	if err != nil {
		return err
	}
	return writeTo(trie, opt.out)
}

func dumpDarts(file io.Reader) error {
	ss, err := asSlice(file)
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

func fromLines(file io.Reader) (*keytree.Tree, error) {
	ks := keytree.New()
	var i uint32
	t := token.NewLinedWords(file)
	err := t.Walk(func(w word.Word) error {
		defer func() { i++ }()
		if err := ks.Put(w, i); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ks, nil
}

func asSlice(r io.Reader) ([]string, error) {
	var ss []string
	err := token.NewLinedString(r).Walk(func(s string) error {
		ss = append(ss, s)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ss, nil
}
