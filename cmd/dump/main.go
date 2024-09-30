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
	"github.com/ajiyoshi-vg/hairetsu/codec/bytes"
	"github.com/ajiyoshi-vg/hairetsu/codec/composer"
	"github.com/ajiyoshi-vg/hairetsu/codec/runes"
	"github.com/ajiyoshi-vg/hairetsu/codec/u16s"
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
	case "darts":
		return dumpDarts(file)
	case "bytes-m":
		return composeBytes(file, bytes.NewMapDict())
	case "bytes-i":
		return composeBytes(file, bytes.NewIdentityDict())
	case "bytes-a":
		return composeBytes(file, bytes.NewArrayDict())
	case "u16s-m":
		return composeU16s(file, u16s.NewMapDict())
	case "u16s-i":
		return composeU16s(file, u16s.NewIdentityDict())
	case "u16s-a":
		return composeU16s(file, u16s.NewArrayDict())
	case "runes-m":
		return composeRunes(file, runes.NewMapDict())
	case "runes-i":
		return composeRunes(file, runes.NewIdentityDict())
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

func composeBytes[D bytes.WordDict](r io.ReadSeeker, dict D) error {
	c := composer.NewBytes(dict, options()...)
	trie, err := c.Compose(r)
	if err != nil {
		return err
	}
	return writeTo(trie, opt.out)
}
func composeU16s[D u16s.WordDict](r io.ReadSeeker, dict D) error {
	c := composer.NewInt16(dict, options()...)
	trie, err := c.Compose(r)
	if err != nil {
		return err
	}
	return writeTo(trie, opt.out)
}
func composeRunes[D runes.WordDict](r io.ReadSeeker, dict D) error {
	c := composer.NewRunes(dict, options()...)
	trie, err := c.Compose(r)
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
