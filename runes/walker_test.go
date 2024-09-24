package runes

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"testing"

	"github.com/ajiyoshi-vg/external/scan"
	"github.com/ajiyoshi-vg/hairetsu/doublearray/item"
	"github.com/ajiyoshi-vg/hairetsu/word"
	"github.com/stretchr/testify/assert"
)

type mock struct {
	item []item.Item
}

func (m *mock) Put(i item.Item) {
	m.item = append(m.item, i)
}

func (m *mock) Get(w word.Word) error {
	for _, i := range m.item {
		if word.Equal(i.Word, w) {
			return nil
		}
	}
	return fmt.Errorf("not found")
}

func TestStringLines(t *testing.T) {
	cases := []struct {
		title  string
		input  io.Reader
		expect []string
	}{
		{
			title: "normal",
			input: bytes.NewBufferString("aaa\nbb\ncccc"),
			expect: []string{
				"aaa",
				"bb",
				"cccc",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			ss := slices.Collect(scan.Lines(c.input))
			mock := &mock{}
			dict, err := FromSlice(ss, mock)
			assert.NoError(t, err)
			for _, s := range c.expect {
				w, err := dict.Word(s)
				assert.NoError(t, err)
				assert.NoError(t, mock.Get(w))
			}
		})
	}
}
