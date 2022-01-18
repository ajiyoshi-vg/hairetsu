package token

import (
	"bufio"
	"bytes"
	"io"

	"github.com/ajiyoshi-vg/hairetsu/word"
)

type LinedString struct {
	r io.Reader
}
type LinedBytes struct {
	r io.Reader
}
type LinedWords struct {
	r io.Reader
}

func NewLinedString(r io.Reader) *LinedString {
	return &LinedString{r: r}
}

func NewLinedBytes(r io.Reader) *LinedBytes {
	return &LinedBytes{r: r}
}
func NewLinedWords(r io.Reader) *LinedWords {
	return &LinedWords{r: r}
}

func (x *LinedString) Walk(yield func(s string) error) error {
	tee := &bytes.Buffer{}
	scan := bufio.NewScanner(io.TeeReader(x.r, tee))
	x.r = tee
	for i := 0; scan.Scan(); i++ {
		line := scan.Text()
		if err := yield(line); err != nil {
			return err
		}
		if err := scan.Err(); err != nil {
			return err
		}
	}
	return nil
}

func (x *LinedBytes) Walk(yield func(s []byte) error) error {
	tee := &bytes.Buffer{}
	scan := bufio.NewScanner(io.TeeReader(x.r, tee))
	x.r = tee
	for i := 0; scan.Scan(); i++ {
		line := scan.Bytes()
		if err := yield(line); err != nil {
			return err
		}
		if err := scan.Err(); err != nil {
			return err
		}
	}
	return nil
}

func (x *LinedWords) Walk(yield func(s word.Word) error) error {
	tee := &bytes.Buffer{}
	scan := bufio.NewScanner(io.TeeReader(x.r, tee))
	x.r = tee
	for i := 0; scan.Scan(); i++ {
		line := scan.Bytes()
		if err := yield(word.FromBytes(line)); err != nil {
			return err
		}
		if err := scan.Err(); err != nil {
			return err
		}
	}
	return nil
}
