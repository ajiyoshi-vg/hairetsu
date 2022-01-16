package lines

import (
	"bufio"
	"io"

	"github.com/ajiyoshi-vg/hairetsu/word"
)

func StringSlice(r io.Reader) ([]string, error) {
	ret := make([]string, 0, 10)
	err := AsString(r, func(s string) error {
		ret = append(ret, s)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func AsString(r io.Reader, yield func(string) error) error {
	scan := bufio.NewScanner(r)
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

func AsWords(r io.Reader, yield func(word.Word) error) error {
	scan := bufio.NewScanner(r)
	for i := 0; scan.Scan(); i++ {
		line := scan.Text()
		if err := yield(word.FromBytes([]byte(line))); err != nil {
			return err
		}
		if err := scan.Err(); err != nil {
			return err
		}
	}
	return nil
}
