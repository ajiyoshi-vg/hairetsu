package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"time"

	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/progress"
	"github.com/ikawaha/dartsclone"
	dartsprog "github.com/ikawaha/dartsclone/progressbar"
)

type option struct {
	in      string
	out     string
	kind    string
	verbose bool
}

var opt option

func init() {
	flag.StringVar(&opt.in, "in", "bench.dat", "line sep text default: bench.txt")
	flag.StringVar(&opt.out, "o", "out.dat", "output")
	flag.StringVar(&opt.kind, "kind", "byte", "[rune|byte|dict|darts] default: byte")
	flag.BoolVar(&opt.verbose, "v", false, "verbose")
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

func options() []doublearray.Option {
	ret := []doublearray.Option{
		doublearray.OptionProgress(&progress.ProgressBar{}),
	}
	if opt.verbose {
		return append(ret, doublearray.Verbose)
	}
	return ret
}

func dumpByte(file io.Reader) error {
	trie, err := hairetsu.NewByteTrieBuilder(options()...).BuildFromLines(file)
	if err != nil {
		return err
	}
	return writeTo(trie, opt.out)
}

func dumpRune(file io.ReadSeeker) error {
	trie, err := hairetsu.NewRuneTrieBuilder(options()...).BuildFromLines(file)
	if err != nil {
		return err
	}
	return writeTo(trie, opt.out)
}

func dumpDict(file io.ReadSeeker) error {
	trie, err := hairetsu.NewDictTrieBuilder(options()...).BuildFromLines(file)
	if err != nil {
		return err
	}
	return writeTo(trie, opt.out)
}

func dumpDarts(file io.Reader) error {
	ss := slices.Collect(scan.Lines(file))

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
