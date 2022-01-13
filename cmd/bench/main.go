package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ajiyoshi-vg/hairetsu/doublearray"
	da "github.com/ajiyoshi-vg/hairetsu/doublearray"
	"github.com/ajiyoshi-vg/hairetsu/keyset"
	"github.com/ajiyoshi-vg/hairetsu/keytree"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/ikawaha/dartsclone"
	progress "github.com/ikawaha/dartsclone/progressbar"
	"github.com/pkg/profile"
	"github.com/schollz/progressbar"
)

type option struct {
	size int
	kind string
}

var opt option

func init() {
	flag.IntVar(&opt.size, "size", 1000*1000, "# of key")
	flag.StringVar(&opt.kind, "kind", "tree", "[set|tree|darts] default:tree")
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
	err := benchmark()
	fmt.Println("finished")
	return err
}

func benchmark() error {
	switch opt.kind {
	case "set":
		return benchKeySet(opt.size)
	case "tree":
		return benchKeyTree(opt.size)
	case "darts":
		return benchDartsClose(opt.size)
	default:
		return fmt.Errorf("uknown kind %s", opt.kind)
	}
}

func benchKeySet(num int) error {
	data := keysetData(num)
	return benchHairetsu(data)
}

func benchKeyTree(num int) error {
	data := keytreeData(num)
	return benchHairetsu(data)
}

func benchDartsClose(num int) error {
	data := dartsCloseData(num)
	pbar := progress.New()
	pbar.SetMaximum(num)
	t, err := dartsclone.BuildTRIE(data, nil, pbar)
	if err != nil {
		return err
	}
	for i, x := range data {
		actual, _, err := t.ExactMatchSearch(x)
		if err != nil {
			return err
		}
		if i != actual {
			return fmt.Errorf("(%v) want %d got %d", []byte(x), i, actual)
		}
	}
	return nil
}

func dartsCloseData(num int) []string {
	keys := make([]string, 0, num)
	for i := 1; i < num; i++ {
		keys = append(keys, fmt.Sprintf("%d", i))
	}
	return keys
}

func keysetData(num int) keyset.KeySet {
	data := make(keyset.KeySet, 0, num)
	for i := 0; i < num; i++ {
		buf := []byte(fmt.Sprintf("%d", i))
		data = append(data, keyset.Item{
			Key: word.FromBytes(buf),
			Val: uint32(i),
		})
	}
	return data
}

func keytreeData(num int) *keytree.Tree {
	data := keytree.New()
	for i := 0; i < num; i++ {
		buf := []byte(fmt.Sprintf("%d", i))
		data.Put(word.FromBytes(buf), uint32(i))
	}
	return data
}

func benchHairetsu(data doublearray.Walker) error {
	x := da.New(10)
	if err := da.NewBuilder(da.OptionProgress(progressbar.New(1))).Build(x, data); err != nil {
		return err
	}
	log.Println("build finished")
	log.Println(x.Stat())

	return nil
}
