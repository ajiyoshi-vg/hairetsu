package composer

import (
	"io"
	"os"

	"github.com/ajiyoshi-vg/hairetsu/codec/u16s"
	"github.com/ajiyoshi-vg/hairetsu/doublearray"
)

type FileInline[D u16s.WordDict] struct {
	dict D
	data *doublearray.DoubleArray
}

func NewFileInline[D u16s.WordDict](dict D) *FileInline[D] {
	return &FileInline[D]{dict: dict}
}

func (f *FileInline[D]) ReadFrom(r io.Reader) (int64, error) {
	var ret int64
	n, err := f.dict.ReadFrom(r)
	ret += n
	if err != nil {
		return ret, err
	}
	da := doublearray.New()
	m, err := da.ReadFrom(r)
	ret += m
	if err != nil {
		return ret, err
	}
	f.data = da
	return ret, nil
}

func (f *FileInline[D]) Open(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := f.ReadFrom(file); err != nil {
		return err
	}
	return nil
}

func (f *FileInline[D]) Searcher() *u16s.InlineSearcher[D, *doublearray.DoubleArray] {
	return u16s.NewInlineSearcher(f.data, f.dict)
}
