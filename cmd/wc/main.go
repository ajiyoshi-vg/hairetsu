package main

import (
	"cmp"
	"flag"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"maps"
	"os"
	"slices"

	"github.com/ajiyoshi-vg/hairetsu"
	"github.com/ajiyoshi-vg/hairetsu/codec/u16s"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

var opt struct {
	in   string
	mode string
}

func init() {
	flag.StringVar(&opt.in, "in", "", "input file")
	flag.StringVar(&opt.mode, "mode", "byte", "byte | uint16")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
	}
}
func run() error {
	file, err := os.Open(opt.in)
	if err != nil {
		return err
	}
	defer file.Close()
	return process(file)
}

func process(r io.Reader) error {
	switch opt.mode {
	case "byte":
		da := doublearray.New()
		if _, err := da.ReadFrom(r); err != nil {
			return err
		}
		count(byteFromWord(wordFromItem(doublearray.Leafs(da))))
	case "uint16":
		da := doublearray.New()
		if _, err := da.ReadFrom(r); err != nil {
			return err
		}
		count(uint16FromByte(bytesFromWord(wordFromItem(doublearray.Leafs(da)))))
	case "rune":
		da := &hairetsu.RuneTrie{}
		if _, err := da.ReadFrom(r); err != nil {
			return err
		}
		count(runeFromWord(wordFromItem(da.Leafs())))
	case "double":
		da := hairetsu.NewDoubleByteTrie(nil, u16s.NewMapDict())
		if _, err := da.ReadFrom(r); err != nil {
			return err
		}
		count(runeFromWord(wordFromItem(da.Leafs())))
	}
	return nil
}

func count[T cmp.Ordered](seq iter.Seq[T]) {
	wc := make(map[T]int)
	total := 0
	for x := range seq {
		wc[x]++
		total++
	}
	n := 0
	keys := slices.Collect(maps.Keys(wc))
	slices.Sort(keys)
	slices.Reverse(keys)
	for _, key := range keys {
		c := wc[key]
		rate := fmt.Sprintf("%.2f%%", 100*float64(c)/float64(total))
		slog.Info("wc", "word", key, "count", c, "rate", rate)
		n++
	}
	slog.Info("wc", "total", int64(total), "unique", n)
}

func wordFromItem(seq iter.Seq[item.Item]) iter.Seq[word.Word] {
	return func(yield func(word.Word) bool) {
		for x := range seq {
			if !yield(x.Word) {
				return
			}
		}
	}
}

func runeFromWord(seq iter.Seq[word.Word]) iter.Seq[rune] {
	return func(yield func(rune) bool) {
		for x := range seq {
			for _, c := range x {
				if !yield(rune(c)) {
					return
				}
			}
		}
	}
}

func byteFromWord(seq iter.Seq[word.Word]) iter.Seq[byte] {
	return func(yield func(byte) bool) {
		for x := range seq {
			bs, err := x.Bytes()
			if err != nil {
				slog.Error(err.Error())
				return
			}
			for _, b := range bs {
				if !yield(b) {
					return
				}
			}
		}
	}
}
func bytesFromWord(seq iter.Seq[word.Word]) iter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		for x := range seq {
			bs, err := x.Bytes()
			if err != nil {
				slog.Error(err.Error())
				return
			}
			if !yield(bs) {
				return
			}
		}
	}
}

func uint16FromByte(seq iter.Seq[[]byte]) iter.Seq[uint16] {
	return func(yield func(uint16) bool) {
		for bs := range seq {
			for i := 0; i < len(bs); i += 2 {
				var val uint16
				for j := 0; j < 2; j++ {
					if i+j < len(bs) {
						val |= uint16(bs[i+j]) << (8 * uint(j))
					}
				}
				if !yield(val) {
					return
				}
			}
		}
	}
}
