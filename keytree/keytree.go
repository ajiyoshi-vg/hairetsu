package keytree

import (
	"bytes"
	"fmt"
	"io"

	"github.com/ajiyoshi-vg/hairetsu/lines"
	"github.com/ajiyoshi-vg/hairetsu/runes"
	"github.com/ajiyoshi-vg/hairetsu/word"
)

type Tree struct {
	value    *uint32
	children map[word.Code]*Tree
	leafNum  int
}

func New() *Tree {
	return &Tree{
		children: map[word.Code]*Tree{},
	}
}

func FromBytes(xs [][]byte) (*Tree, error) {
	ret := New()
	for i, x := range xs {
		err := ret.Put(word.FromBytes(x), uint32(i))
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func FromLines(file io.Reader) (*Tree, error) {
	ks := New()
	var i uint32
	err := lines.AsWords(file, func(w word.Word) error {
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

func FromStringLines(r io.Reader) (*Tree, runes.Dict, error) {
	tee := &bytes.Buffer{}
	dict, err := runes.FromLines(io.TeeReader(r, tee))
	if err != nil {
		return nil, nil, err
	}

	ks := New()
	var i uint32
	err = lines.AsString(tee, func(s string) error {
		defer func() { i++ }()
		w, err := dict.Word(s)
		if err != nil {
			return err
		}
		return ks.Put(w, i)
	})
	if err != nil {
		return nil, nil, err
	}
	return ks, dict, err
}

func FromWord(data []word.Word) *Tree {
	root := New()
	for i, x := range data {
		root.Put(x, uint32(i))
	}
	return root
}

func (x *Tree) LeafNum() int {
	return x.leafNum
}

func (x *Tree) Get(key word.Word) (*uint32, error) {
	node := x
	for _, b := range key {
		node = node.children[b]
		if node == nil {
			return nil, fmt.Errorf("not found")
		}
	}
	return node.value, nil
}

func (x *Tree) Put(key word.Word, val uint32) error {
	node := x
	for _, b := range key {
		child := node.children[b]
		if child == nil {
			child = New()
			node.children[b] = child
		}
		node = child
	}
	if node.value != nil {
		return fmt.Errorf("%v was inserted twice. old:%d new:%d", key, *node.value, val)
	}
	node.value = &val
	x.leafNum++
	return nil
}

func (x *Tree) WalkNode(f func(word.Word, []word.Code, *uint32) error) error {
	return x.walkNode(word.Word{}, f)
}

func (x *Tree) WalkLeaf(f func(word.Word, uint32) error) error {
	return x.walkLeaf(word.Word{}, f)
}

func (x *Tree) walkNode(prefix word.Word, f func(word.Word, []word.Code, *uint32) error) error {
	branch := make([]word.Code, 0, len(x.children))
	for b := range x.children {
		branch = append(branch, b)
	}
	if err := f(prefix, branch, x.value); err != nil {
		return err
	}
	for b, child := range x.children {
		if err := child.walkNode(append(prefix, b), f); err != nil {
			return err
		}
	}
	return nil
}

func (x *Tree) walkLeaf(prefix word.Word, f func(word.Word, uint32) error) error {
	if x.value != nil {
		if err := f(prefix, *x.value); err != nil {
			return err
		}
	}
	for b, child := range x.children {
		if err := child.walkLeaf(append(prefix, b), f); err != nil {
			return err
		}
	}
	return nil
}
