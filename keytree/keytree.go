package keytree

import (
	"errors"
	"fmt"

	"github.com/ajiyoshi-vg/hairetsu/word"
)

// Deprecated: use doublearray.Builder.Factory() instead
type Tree struct {
	value    *uint32
	children map[word.Code]*Tree
	leafNum  int
}

// Deprecated: use doublearray.Builder.Factory() instead
func New() *Tree {
	return &Tree{
		children: map[word.Code]*Tree{},
	}
}

// Deprecated: use doublearray.Builder.Factory() instead
func FromWord(data []word.Word) *Tree {
	root := New()
	for i, x := range data {
		_ = root.Put(x, uint32(i))
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

var ErrDuplicated = errors.New("key already exists")

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
		return ErrDuplicated
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
